package matcher

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/CameronHonis/set"
	"math"
	"sync"
	"time"
)

// NOTE: generator compilation broken
// //go:generate mockgen -destination matcher_service_mock.go . MatcherServiceI
type MatcherServiceI interface {
	ServiceI
	MatchById(matchId string) (*models.Match, error)
	MatchByClientKey(clientKey models.Key) (*models.Match, error)
	InboundChallenges(challengedKey models.Key) (*Set[*models.Challenge], error)
	OutboundChallenges(challengerKey models.Key) (*Set[*models.Challenge], error)
	GetChallenge(challengerKey, receivingClientKey models.Key) (*models.Challenge, error)

	ExecuteMove(matchId string, move *chess.Move) error

	RequestChallenge(challenge *models.Challenge) error
	AcceptChallenge(challengedKey, challengerKey models.Key) error
	RevokeChallenge(challengerKey, challengedKey models.Key) error
	DeclineChallenge(challengedKey, challengerKey models.Key) error
	AddMatch(match *models.Match) error
}

type MatcherService struct {
	Service
	__dependencies__ Marker
	LogService       LoggerServiceI
	AuthService      auth.AuthenticationServiceI

	__state__            Marker
	matchByMatchId       map[string]*models.Match
	matchIdByClientKey   map[models.Key]string
	outboundsByClientKey map[models.Key]*Set[*models.Challenge]
	inboundsByClientKey  map[models.Key]*Set[*models.Challenge]
	mu                   sync.Mutex
}

func NewMatcherService(config *MatcherServiceConfig) *MatcherService {
	matchService := &MatcherService{
		matchByMatchId:       make(map[string]*models.Match),
		matchIdByClientKey:   make(map[models.Key]string),
		outboundsByClientKey: make(map[models.Key]*Set[*models.Challenge]),
		inboundsByClientKey:  make(map[models.Key]*Set[*models.Challenge]),
	}
	matchService.Service = *NewService(matchService, config)
	return matchService
}

func (m *MatcherService) MatchById(matchId string) (*models.Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	match, ok := m.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("matcher with id %s not found", matchId)
	}
	return match, nil
}

func (m *MatcherService) MatchByClientKey(clientKey models.Key) (*models.Match, error) {
	m.mu.Lock()
	matchId, ok := m.matchIdByClientKey[clientKey]
	m.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("client %s not in matcher", clientKey)
	}
	return m.MatchById(matchId)
}

func (m *MatcherService) InboundChallenges(challengedKey models.Key) (*Set[*models.Challenge], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	challenges, ok := m.inboundsByClientKey[challengedKey]
	if !ok {
		challenges = EmptySet[*models.Challenge]()
		m.inboundsByClientKey[challengedKey] = challenges
	}
	return challenges, nil
}

func (m *MatcherService) OutboundChallenges(challengerKey models.Key) (*Set[*models.Challenge], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	challenges, ok := m.outboundsByClientKey[challengerKey]
	if !ok {
		challenges = EmptySet[*models.Challenge]()
		m.outboundsByClientKey[challengerKey] = challenges
	}
	return challenges, nil
}

func (m *MatcherService) GetChallenge(challengerKey models.Key, receiverKey models.Key) (*models.Challenge, error) {
	outbounds, outboundsErr := m.OutboundChallenges(challengerKey)
	if outboundsErr != nil {
		return nil, outboundsErr
	}

	inbounds, inboundsErr := m.InboundChallenges(receiverKey)
	if inboundsErr != nil {
		return nil, inboundsErr
	}

	matches := outbounds.Intersect(inbounds).Flatten()
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find challenge from %s to %s", challengerKey, receiverKey)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("more than one challenge from %s to %s", challengerKey, receiverKey)
	}

	return matches[0], nil
}

func (m *MatcherService) ExecuteMove(matchId string, move *chess.Move) error {
	match, getMatchErr := m.MatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	if !chess.IsLegalMove(match.Board, move) {
		return fmt.Errorf("move %v is not legal", move)
	}

	matchBuilder := models.NewMatchBuilder().FromMatch(match)
	currTime := time.Now()
	matchBuilder.WithLastMoveTime(&currTime)
	secondsSinceLastMove := math.Max(currTime.Sub(*match.LastMoveTime).Seconds(), 0.1)
	if match.Board.IsWhiteTurn {
		newWhiteTimeRemaining := match.WhiteTimeRemainingSec - math.Max(0.1, secondsSinceLastMove)
		matchBuilder.WithWhiteTimeRemainingSec(math.Max(0, newWhiteTimeRemaining))
		if newWhiteTimeRemaining == 0 {
			boardBuilder := chess.NewBoardBuilder().FromBoard(match.Board)
			boardBuilder.WithIsTerminal(true)
			boardBuilder.WithIsBlackWinner(true)
			matchBuilder.WithBoard(boardBuilder.Build())
		}
	} else {
		newBlackTimeRemaining := match.BlackTimeRemainingSec - math.Max(0.1, secondsSinceLastMove)
		matchBuilder.WithBlackTimeRemainingSec(math.Max(0, newBlackTimeRemaining))
		if newBlackTimeRemaining == 0 {
			boardBuilder := chess.NewBoardBuilder().FromBoard(match.Board)
			boardBuilder.WithIsTerminal(true)
			boardBuilder.WithIsWhiteWinner(true)
			matchBuilder.WithBoard(boardBuilder.Build())
		}
	}
	newBoard := chess.GetBoardFromMove(match.Board, move)
	matchBuilder.WithBoard(newBoard)
	newMatch := matchBuilder.Build()

	setMatchErr := m.SetMatch(newMatch)
	if setMatchErr != nil {
		return setMatchErr
	}

	if newMatch.Board.IsTerminal {
		return m.RemoveMatch(newMatch)
	} else {
		go m.StartTimer(newMatch)
	}

	return nil
}

func (m *MatcherService) RequestChallenge(challenge *models.Challenge) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("client %s challenging client %s", challenge.ChallengerKey, challenge.ChallengedKey))
	if challengeErr := m.ValidateChallenge(challenge); challengeErr != nil {
		go m.Dispatch(NewChallengeRequestFailedEvent(challenge, challengeErr.Error()))
		return challengeErr
	}

	isBotChallenge := challenge.BotName != ""
	if isBotChallenge {
		// NOTE: bot challenge needs a bot server client key
		challenge.ChallengedKey = m.AuthService.ClientKeysByRole(models.BOT).Flatten()[0]
	}
	challengerOutbounds, _ := m.OutboundChallenges(challenge.ChallengerKey)
	challengedInbounds, _ := m.InboundChallenges(challenge.ChallengedKey)

	m.mu.Lock()
	defer m.mu.Unlock()
	challengerOutbounds.Add(challenge)
	challengedInbounds.Add(challenge)

	go m.Dispatch(NewChallengeCreatedEvent(challenge))
	return nil
}

func (m *MatcherService) AcceptChallenge(challengedKey, challengerKey models.Key) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("accepting challenge with client %s", challengedKey))
	challenge, challengeErr := m.GetChallenge(challengerKey, challengedKey)
	if challengeErr != nil {
		go m.Dispatch(NewMatchCreationFailedEvent(challengerKey, "challenged unavailable for matcher"))
		return challengeErr
	}
	match := models.NewMatchBuilder().FromChallenge(challenge).Build()
	return m.AddMatch(match)
}

func (m *MatcherService) RevokeChallenge(challengerKey, challengedKey models.Key) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("canceling challenge for challenger %s", challengerKey))
	panic("implement me")
	return nil
}

func (m *MatcherService) DeclineChallenge(challengedKey, challengerKey models.Key) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("revoking challenge for challenger %s", challengerKey))
	panic("implement me")
	return nil
}

func (m *MatcherService) AddMatch(match *models.Match) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("adding matcher %s", match.Uuid))
	if !m.CanStartMatchWithClientKey(match.WhiteClientKey) {
		go m.Dispatch(NewMatchCreationFailedEvent(match.WhiteClientKey, "white client unavailable for matcher"))
		return fmt.Errorf("white client %s unavailable for matcher", match.WhiteClientKey)
	}
	if !m.CanStartMatchWithClientKey(match.BlackClientKey) {
		go m.Dispatch(NewMatchCreationFailedEvent(match.BlackClientKey, "black client unavailable for matcher"))
		return fmt.Errorf("black client %s unavailable for matcher", match.BlackClientKey)
	}

	m.mu.Lock()
	m.matchByMatchId[match.Uuid] = match
	if role, _ := m.AuthService.GetRole(match.WhiteClientKey); role != models.BOT {
		m.matchIdByClientKey[match.WhiteClientKey] = match.Uuid
	}
	if role, _ := m.AuthService.GetRole(match.BlackClientKey); role != models.BOT {
		m.matchIdByClientKey[match.BlackClientKey] = match.Uuid
	}
	m.mu.Unlock()

	go m.Dispatch(NewMatchCreatedEvent(match))
	return nil
}

func (m *MatcherService) SetMatch(newMatch *models.Match) error {
	oldMatch, fetchCurrMatchErr := m.MatchById(newMatch.Uuid)
	if fetchCurrMatchErr != nil {
		return fetchCurrMatchErr
	}
	if newMatch.WhiteClientKey != oldMatch.WhiteClientKey {
		return fmt.Errorf("cannot change white client id")
	}
	if newMatch.BlackClientKey != oldMatch.BlackClientKey {
		return fmt.Errorf("cannot change black client id")
	}
	if !newMatch.TimeControl.Equals(oldMatch.TimeControl) {
		return fmt.Errorf("cannot change time control")
	}
	m.mu.Lock()
	m.matchByMatchId[newMatch.Uuid] = newMatch
	m.mu.Unlock()

	go m.Dispatch(NewMatchUpdated(newMatch))
	return nil
}

func (m *MatcherService) RemoveMatch(match *models.Match) error {
	m.LogService.Log(models.ENV_MATCH_SERVICE, fmt.Sprintf("removing matcher %s", match.Uuid))
	m.mu.Lock()
	if _, ok := m.matchByMatchId[match.Uuid]; !ok {
		return fmt.Errorf("matcher with id %s doesn't exist", match.Uuid)
	}
	if match.WhiteClientKey != "" {
		delete(m.matchIdByClientKey, match.WhiteClientKey)
	}
	if match.BlackClientKey != "" {
		delete(m.matchIdByClientKey, match.BlackClientKey)
	}
	delete(m.matchByMatchId, match.Uuid)
	m.mu.Unlock()

	go m.Dispatch(NewMatchEndedEvent(match))
	return nil
}

func (m *MatcherService) CanStartMatchWithClientKey(clientKey models.Key) bool {
	role, roleErr := m.AuthService.GetRole(clientKey)
	if roleErr != nil {
		return false
	}
	if role == models.BOT {
		return true
	}

	match, _ := m.MatchByClientKey(clientKey)
	return match == nil
}

func (m *MatcherService) ValidateChallenge(challenge *models.Challenge) error {
	if challenge.ChallengedKey == "" {
		if challenge.BotName == "" {
			return fmt.Errorf("challenged key and bot name cannot both be zero values")
		}
		if !m.AuthService.BotClientExists() {
			return fmt.Errorf("bot server offline")
		}
	} else {
		if challenge.BotName != "" {
			return fmt.Errorf("challenged key and bot name cannot both be populated")
		}
	}
	if challenge.ChallengerKey == challenge.ChallengedKey {
		return fmt.Errorf("cannot challenge self")
	}
	if !m.CanStartMatchWithClientKey(challenge.ChallengerKey) {
		return fmt.Errorf("challenger %s unavailable for matcher", challenge.ChallengerKey)
	}
	if challengeDuplicate, _ := m.GetChallenge(challenge.ChallengerKey, challenge.ChallengedKey); challengeDuplicate != nil {
		return fmt.Errorf("challenge already exists")
	}
	return nil
}

func (m *MatcherService) StartTimer(match *models.Match) {
	var waitTime time.Duration
	if match.Board.IsWhiteTurn {
		waitTime = time.Duration(match.WhiteTimeRemainingSec) * time.Second
	} else {
		waitTime = time.Duration(match.BlackTimeRemainingSec) * time.Second
	}

	time.Sleep(waitTime)
	currMatch, _ := m.MatchById(match.Uuid)
	if currMatch == nil {
		m.LogService.LogRed(models.ENV_TIMER, "matcher not found")
		return
	}
	if currMatch.LastMoveTime.Equal(*match.LastMoveTime) {
		matchBuilder := models.NewMatchBuilder().FromMatch(match)
		boardBuilder := chess.NewBoardBuilder().FromBoard(match.Board)
		boardBuilder.WithIsTerminal(true)
		if match.Board.IsWhiteTurn {
			matchBuilder.WithWhiteTimeRemainingSec(0)
			boardBuilder.WithIsBlackWinner(true)
		} else {
			matchBuilder.WithBlackTimeRemainingSec(0)
			boardBuilder.WithIsWhiteWinner(true)
		}
		matchBuilder.WithBoard(boardBuilder.Build())
		newMatch := matchBuilder.Build()
		_ = m.SetMatch(newMatch)
	}
}
