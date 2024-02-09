package matcher

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/sub_service"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/CameronHonis/set"
	"math"
	"sync"
	"time"
)

type MatcherServiceI interface {
	service.ServiceI
	MatchById(matchId string) (*models.Match, error)
	MatchByClientKey(clientKey models.Key) (*models.Match, error)
	InboundChallenges(challengedKey models.Key) (*set.Set[*models.Challenge], error)
	OutboundChallenges(challengerKey models.Key) (*set.Set[*models.Challenge], error)
	GetChallenge(challengerKey, receivingClientKey models.Key) (*models.Challenge, error)

	ExecuteMove(matchId string, move *chess.Move) error
	ResignMatch(matchId string, clientKey models.Key) error

	RequestChallenge(challenge *models.Challenge) error
	AcceptChallenge(challengedKey, challengerKey models.Key) error
	RevokeChallenge(challengerKey, challengedKey models.Key) error
	DeclineChallenge(challengerKey, challengedKey models.Key) error

	AddMatch(match *models.Match) error
}

type MatcherService struct {
	service.Service
	__dependencies__ marker.Marker
	Logger           log.LoggerServiceI
	AuthService      auth.AuthenticationServiceI
	SubService       sub_service.SubscriptionServiceI

	__state__            marker.Marker
	matchByMatchId       map[string]*models.Match
	matchIdByClientKey   map[models.Key]string
	outboundsByClientKey map[models.Key]*set.Set[*models.Challenge]
	inboundsByClientKey  map[models.Key]*set.Set[*models.Challenge]
	mu                   sync.Mutex
}

func NewMatcherService(config *MatcherServiceConfig) *MatcherService {
	matchService := &MatcherService{
		matchByMatchId:       make(map[string]*models.Match),
		matchIdByClientKey:   make(map[models.Key]string),
		outboundsByClientKey: make(map[models.Key]*set.Set[*models.Challenge]),
		inboundsByClientKey:  make(map[models.Key]*set.Set[*models.Challenge]),
	}
	matchService.Service = *service.NewService(matchService, config)
	return matchService
}

func (m *MatcherService) OnBuild() {
	m.AddEventListener(MATCH_UPDATED, OnMatchUpdated)
}

func (m *MatcherService) MatchById(matchId string) (*models.Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	match, ok := m.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not found", matchId)
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

func (m *MatcherService) InboundChallenges(challengedKey models.Key) (*set.Set[*models.Challenge], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	challenges, ok := m.inboundsByClientKey[challengedKey]
	if !ok {
		challenges = set.EmptySet[*models.Challenge]()
		m.inboundsByClientKey[challengedKey] = challenges
	}
	return challenges, nil
}

func (m *MatcherService) OutboundChallenges(challengerKey models.Key) (*set.Set[*models.Challenge], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	challenges, ok := m.outboundsByClientKey[challengerKey]
	if !ok {
		challenges = set.EmptySet[*models.Challenge]()
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
	m.Logger.Log(models.ENV_MATCHER_SERVICE, "executing move on match ", matchId)
	match, getMatchErr := m.MatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	if !chess.IsLegalMove(match.Board, move) {
		return fmt.Errorf("move %v is not legal", move)
	}

	matchBuilder := builders.NewMatchBuilder().FromMatch(match)
	currTime := time.Now()
	matchBuilder.WithLastMoveTime(&currTime)
	secondsSinceLastMove := math.Max(currTime.Sub(*match.LastMoveTime).Seconds(), 0.1)
	if match.Board.IsWhiteTurn {
		newWhiteTimeRemaining := match.WhiteTimeRemainingSec - math.Max(0.1, secondsSinceLastMove)
		matchBuilder.WithWhiteTimeRemainingSec(math.Max(0, newWhiteTimeRemaining))
		if newWhiteTimeRemaining == 0 {
			matchBuilder.WithResult(models.MATCH_RESULT_BLACK_WINS_BY_TIMEOUT)
		}
	} else {
		newBlackTimeRemaining := match.BlackTimeRemainingSec - math.Max(0.1, secondsSinceLastMove)
		matchBuilder.WithBlackTimeRemainingSec(math.Max(0, newBlackTimeRemaining))
		if newBlackTimeRemaining == 0 {
			matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT)
		}
	}
	newBoard := chess.GetBoardFromMove(match.Board, move)
	matchBuilder.WithBoard(newBoard)
	newMatch := matchBuilder.Build()

	setMatchErr := m.SetMatch(newMatch)
	if setMatchErr != nil {
		return setMatchErr
	}

	return nil
}

func (m *MatcherService) ResignMatch(matchId string, clientKey models.Key) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("client %s resigning match %s", clientKey, matchId))
	match, matchErr := m.MatchById(matchId)
	if matchErr != nil {
		return matchErr
	}
	var result models.MatchResult
	if match.Board.IsWhiteTurn {
		result = models.MATCH_RESULT_BLACK_WINS_BY_RESIGNATION
	} else {
		result = models.MATCH_RESULT_WHITE_WINS_BY_RESIGNATION
	}

	matchBuilder := builders.NewMatchBuilder().FromMatch(match)
	matchBuilder.WithResult(result)
	newMatch := matchBuilder.Build()

	if setMatchErr := m.SetMatch(newMatch); setMatchErr != nil {
		return setMatchErr
	}
	return nil
}

func (m *MatcherService) RequestChallenge(_challenge *models.Challenge) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("client %s challenging client %s", _challenge.ChallengerKey, _challenge.ChallengedKey))

	now := time.Now()
	challengeBuilder := builders.NewChallengeBuilder()
	challengeBuilder.FromChallenge(_challenge)
	challengeBuilder.WithRandomUuid()
	challengeBuilder.WithIsActive(true)
	challengeBuilder.WithTimeCreated(&now)
	challenge := challengeBuilder.Build()
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

func (m *MatcherService) AcceptChallenge(challengerKey, challengedKey models.Key) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("accepting challenge from client %s to %s", challengerKey, challengedKey))
	challenge, challengeErr := m.GetChallenge(challengerKey, challengedKey)
	if challengeErr != nil {
		go m.Dispatch(NewChallengeAcceptFailedEvent(challenge, "could not get challenge"))
		return challengeErr
	}

	m.mu.Lock()
	m.inboundsByClientKey[challengedKey].Remove(challenge)
	m.outboundsByClientKey[challengerKey].Remove(challenge)
	m.mu.Unlock()

	match := builders.NewMatchBuilder().FromChallenge(challenge).Build()
	if addMatchErr := m.AddMatch(match); addMatchErr != nil {
		return addMatchErr
	}

	go m.Dispatch(NewChallengeAcceptedEvent(challenge))
	return challengeErr
}

func (m *MatcherService) RevokeChallenge(challengerKey, challengedKey models.Key) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("revoking challenge from %s to %s", challengerKey, challengedKey))
	challenge, challengeErr := m.GetChallenge(challengerKey, challengedKey)
	if challengeErr != nil {
		return challengeErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.inboundsByClientKey[challengedKey].Remove(challenge)
	m.outboundsByClientKey[challengerKey].Remove(challenge)

	go m.Dispatch(NewChallengeRevokedEvent(challenge))
	return nil
}

func (m *MatcherService) DeclineChallenge(challengerKey, challengedKey models.Key) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("declining challenge from %s to %s", challengerKey, challengedKey))
	challenge, challengeErr := m.GetChallenge(challengerKey, challengedKey)
	if challengeErr != nil {
		return challengeErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.inboundsByClientKey[challengedKey].Remove(challenge)
	m.outboundsByClientKey[challengerKey].Remove(challenge)

	go m.Dispatch(NewChallengeDeniedEvent(challenge))
	return nil
}

func (m *MatcherService) AddMatch(match *models.Match) error {
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("adding match %s", match.Uuid))
	if whiteAvailableErr := m.validateClientAvailable(match.WhiteClientKey); whiteAvailableErr != nil {
		go m.Dispatch(NewMatchCreationFailedEvent(match.WhiteClientKey, match.BlackClientKey, "white client unavailable for matcher"))
		return fmt.Errorf("white client %s unavailable for matcher: %s", match.WhiteClientKey, whiteAvailableErr.Error())
	}
	if blackAvailableErr := m.validateClientAvailable(match.BlackClientKey); blackAvailableErr != nil {
		go m.Dispatch(NewMatchCreationFailedEvent(match.BlackClientKey, match.BlackClientKey, "black client unavailable for matcher"))
		return fmt.Errorf("black client %s unavailable for matcher: %s", match.BlackClientKey, blackAvailableErr.Error())
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
	m.Logger.Log(models.ENV_MATCHER_SERVICE, fmt.Sprintf("removing match %s", match.Uuid))
	m.mu.Lock()
	if _, ok := m.matchByMatchId[match.Uuid]; !ok {
		return fmt.Errorf("match with id %s doesn't exist", match.Uuid)
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

func (m *MatcherService) validateClientAvailable(clientKey models.Key) error {
	role, roleErr := m.AuthService.GetRole(clientKey)
	if roleErr != nil {
		return fmt.Errorf("could not get role")
	}
	if role == models.BOT {
		// bot clients can manage many matches at a time
		return nil
	}

	if match, _ := m.MatchByClientKey(clientKey); match != nil {
		return fmt.Errorf("client already in match")
	}

	return nil
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

	if challengerAvailableErr := m.validateClientAvailable(challenge.ChallengerKey); challengerAvailableErr != nil {
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
		m.Logger.LogRed(models.ENV_TIMER, "match not found")
		return
	}
	if currMatch.LastMoveTime.Equal(*match.LastMoveTime) {
		matchBuilder := builders.NewMatchBuilder().FromMatch(match)
		if match.Board.IsWhiteTurn {
			matchBuilder.WithWhiteTimeRemainingSec(0)
		} else {
			matchBuilder.WithBlackTimeRemainingSec(0)
		}
		newMatch := matchBuilder.Build()
		_ = m.SetMatch(newMatch)
	}
}

var OnMatchUpdated = func(s service.ServiceI, ev service.EventI) bool {
	matcher := s.(*MatcherService)
	match := ev.Payload().(*MatchUpdatedEventPayload).Match
	if match.Result != models.MATCH_RESULT_IN_PROGRESS {
		_ = matcher.RemoveMatch(match)
	} else {
		go matcher.StartTimer(match)
	}
	return true
}
