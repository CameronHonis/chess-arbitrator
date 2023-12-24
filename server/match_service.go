package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/CameronHonis/set"
	"math"
	"sync"
	"time"
)

type MatchManagerAction string

const (
	MATCH_CREATED     = "MATCH_CREATED"
	MATCH_ENDED       = "MATCH_ENDED"
	MATCH_UPDATED     = "MATCH_UPDATED"
	CHALLENGE_CREATED = "CHALLENGE_CREATED"
	CHALLENGE_DENIED  = "CHALLENGE_DENIED"
	CHALLENGE_REVOKED = "CHALLENGE_REVOKED"
)

type MatchCreatedEventPayload struct {
	Match *Match
}

type MatchCreatedEvent struct{ Event }

func NewMatchCreatedEvent(match *Match) *MatchCreatedEvent {
	return &MatchCreatedEvent{
		Event: *NewEvent(MATCH_CREATED, &MatchCreatedEventPayload{
			Match: match,
		}),
	}
}

type MatchEndedEventPayload struct {
	Match *Match
}

type MatchEndedEvent struct{ Event }

func NewMatchEndedEvent(match *Match) *MatchEndedEvent {
	return &MatchEndedEvent{
		Event: *NewEvent(MATCH_ENDED, &MatchEndedEventPayload{
			Match: match,
		}),
	}
}

type MatchUpdatedEventPayload struct {
	Match *Match
}

type MatchUpdatedEvent struct{ Event }

func NewMatchUpdated(match *Match) *MatchUpdatedEvent {
	return &MatchUpdatedEvent{
		Event: *NewEvent(MATCH_UPDATED, &MatchUpdatedEventPayload{
			Match: match,
		}),
	}
}

type ChallengeCreatedEventPayload struct {
	Challenge *Challenge
}

type ChallengeCreatedEvent struct{ Event }

func NewChallengeCreatedEvent(challenge *Challenge) *ChallengeCreatedEvent {
	return &ChallengeCreatedEvent{
		Event: *NewEvent(CHALLENGE_CREATED, &ChallengeCreatedEventPayload{
			Challenge: challenge,
		}),
	}
}

type ChallengeDeniedEventPayload struct {
	Challenge *Challenge
}

type ChallengeDeniedEvent struct{ Event }

func NewChallengeCanceledEvent(challenge *Challenge) *ChallengeDeniedEvent {
	return &ChallengeDeniedEvent{
		Event: *NewEvent(CHALLENGE_DENIED, &ChallengeDeniedEventPayload{
			Challenge: challenge,
		}),
	}
}

type ChallengeRevokedEventPayload struct {
	Challenge *Challenge
}

type ChallengeRevokedEvent struct{ Event }

func NewChallengeRevokedEvent(challenge *Challenge) *ChallengeRevokedEvent {
	return &ChallengeRevokedEvent{
		Event: *NewEvent(CHALLENGE_REVOKED, &ChallengeRevokedEventPayload{
			Challenge: challenge,
		}),
	}
}

type MatchServiceConfig struct {
	ConfigI
}

func NewMatchServiceConfig() *MatchServiceConfig {
	return &MatchServiceConfig{}
}

type MatchServiceI interface {
	GetMatchById(matchId string) (*Match, error)
	GetMatchByClientKey(clientKey string) (*Match, error)
	GetChallenges(challengerKey string) (*Set[*Challenge], error)
	GetChallenge(challengerKey, receivingClientKey string) (*Challenge, error)

	ExecuteMove(matchId string, move *chess.Move) error

	ChallengePlayer(challenge *Challenge) error
	AcceptChallenge(challengedKey, challengerKey string) error
	RevokeChallenge(challengerKey, challengedKey string) error
	DeclineChallenge(challengedKey, challengerKey string) error
	AddMatch(match *Match) error
}

type MatchService struct {
	Service[*MatchServiceConfig]
	__dependencies__ Marker
	LoggerService    LoggerServiceI
	AuthService      AuthenticationServiceI

	__state__                      Marker
	matchByMatchId                 map[string]*Match
	matchIdByClientKey             map[string]string
	challengeByChallengerClientKey map[string]*Set[*Challenge]
	mu                             sync.Mutex
}

func NewMatchService(config *MatchServiceConfig) *MatchService {
	matchService := &MatchService{
		matchByMatchId:                 make(map[string]*Match),
		matchIdByClientKey:             make(map[string]string),
		challengeByChallengerClientKey: make(map[string]*Set[*Challenge]),
	}
	matchService.Service = *NewService(matchService, config)
	return matchService
}

func (m *MatchService) GetMatchById(matchId string) (*Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	match, ok := m.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not found", matchId)
	}
	return match, nil
}

func (m *MatchService) GetMatchByClientKey(clientKey string) (*Match, error) {
	m.mu.Lock()
	matchId, ok := m.matchIdByClientKey[clientKey]
	m.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("client %s not in match", clientKey)
	}
	return m.GetMatchById(matchId)
}

func (m *MatchService) GetChallenges(challengerKey string) (*Set[*Challenge], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	challenges, ok := m.challengeByChallengerClientKey[challengerKey]
	if !ok {
		return EmptySet[*Challenge](), nil
	}
	return challenges, nil
}

func (m *MatchService) GetChallenge(challengerKey string, receivingClientKey string) (*Challenge, error) {
	challenges, challengesErr := m.GetChallenges(challengerKey)
	if challengesErr != nil {
		return nil, challengesErr
	}
	for _, challenge := range challenges.Flatten() {
		if challenge.ChallengedKey == receivingClientKey {
			return challenge, nil
		}
	}
	return nil, fmt.Errorf("challenge not found")
}

func (m *MatchService) ExecuteMove(matchId string, move *chess.Move) error {
	match, getMatchErr := m.GetMatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	if !chess.IsLegalMove(match.Board, move) {
		return fmt.Errorf("move %v is not legal", move)
	}

	matchBuilder := NewMatchBuilder().FromMatch(match)
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

	setMatchErr := m.setMatch(newMatch)
	if setMatchErr != nil {
		return setMatchErr
	}

	if newMatch.Board.IsTerminal {
		return m.removeMatch(newMatch)
	} else {
		go m.startTimer(newMatch)
	}

	return nil
}

func (m *MatchService) ChallengePlayer(challenge *Challenge) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("client %s challenging client %s", challenge.ChallengerKey, challenge.ChallengedKey))
	if challenge.ChallengerKey == challenge.ChallengedKey {
		return fmt.Errorf("cannot challenge self")
	}
	if !m.canStartMatchWithClientKey(challenge.ChallengerKey) {
		return fmt.Errorf("challenger %s unavailable for match", challenge.ChallengerKey)
	}
	if challengeDuplicate, _ := m.GetChallenge(challenge.ChallengerKey, challenge.ChallengedKey); challengeDuplicate != nil {
		return fmt.Errorf("challenge already exists")
	}
	m.mu.Lock()
	m.challengeByChallengerClientKey[challenge.ChallengerKey].Add(challenge)
	m.mu.Unlock()
	go m.Dispatch(NewChallengeCreatedEvent(challenge))
	return nil
}

func (m *MatchService) AcceptChallenge(challengedKey, challengerKey string) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("accepting challenge with client %s", challengedKey))
	challenge, challengeErr := m.GetChallenge(challengerKey, challengedKey)
	if challengeErr != nil {
		return challengeErr
	}
	match := NewMatchBuilder().FromChallenge(challenge).Build()
	return m.AddMatch(match)
}

func (m *MatchService) RevokeChallenge(challengerKey, challengedKey string) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("canceling challenge for challenger %s", challengerKey))
	panic("implement me")
	return nil
}

func (m *MatchService) DeclineChallenge(challengedKey, challengerKey string) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("revoking challenge for challenger %s", challengerKey))
	panic("implement me")
	return nil
}

func (m *MatchService) AddMatch(match *Match) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("adding match %s", match.Uuid))
	if m.canStartMatchWithClientKey(match.WhiteClientKey) {
		go m.Dispatch(NewMatchCreationFailedEvent(match.WhiteClientKey, "white client unavailable for match"))
		return fmt.Errorf("white client %s unavailable for match", match.WhiteClientKey)
	}
	if m.canStartMatchWithClientKey(match.BlackClientKey) {
		go m.Dispatch(NewMatchCreationFailedEvent(match.BlackClientKey, "black client unavailable for match"))
		return fmt.Errorf("black client %s unavailable for match", match.BlackClientKey)
	}

	m.mu.Lock()
	m.matchByMatchId[match.Uuid] = match
	if role, _ := m.AuthService.GetRole(match.WhiteClientKey); role != BOT {
		m.matchIdByClientKey[match.WhiteClientKey] = match.Uuid
	}
	if role, _ := m.AuthService.GetRole(match.BlackClientKey); role != BOT {
		m.matchIdByClientKey[match.BlackClientKey] = match.Uuid
	}
	m.mu.Unlock()

	go m.Dispatch(NewMatchCreatedEvent(match))
	return nil
}

func (m *MatchService) setMatch(newMatch *Match) error {
	oldMatch, fetchCurrMatchErr := m.GetMatchById(newMatch.Uuid)
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

	m.Dispatch(NewMatchUpdated(newMatch))
	return nil
}

func (m *MatchService) removeMatch(match *Match) error {
	m.LoggerService.Log(ENV_MATCH_SERVICE, fmt.Sprintf("removing match %s", match.Uuid))
	m.mu.Lock()
	defer m.mu.Unlock()
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
	return nil
}

func (m *MatchService) canStartMatchWithClientKey(clientKey string) bool {
	role, roleErr := m.AuthService.GetRole(clientKey)
	if roleErr != nil {
		return false
	}
	if role == BOT {
		return true
	}

	match, _ := m.GetMatchByClientKey(clientKey)
	return match == nil
}

func (m *MatchService) startTimer(match *Match) {
	var waitTime time.Duration
	if match.Board.IsWhiteTurn {
		waitTime = time.Duration(match.WhiteTimeRemainingSec) * time.Second
	} else {
		waitTime = time.Duration(match.BlackTimeRemainingSec) * time.Second
	}

	time.Sleep(waitTime)
	currMatch, _ := m.GetMatchById(match.Uuid)
	if currMatch == nil {
		m.LoggerService.LogRed(ENV_TIMER, "match not found")
		return
	}
	if currMatch.LastMoveTime.Equal(*match.LastMoveTime) {
		matchBuilder := NewMatchBuilder().FromMatch(match)
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
		_ = m.setMatch(newMatch)
	}
}
