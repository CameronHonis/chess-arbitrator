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
	ADD_MATCH           MatchManagerAction = "ADD_MATCH"
	REMOVE_MATCH        MatchManagerAction = "REMOVE_MATCH"
	SET_MATCH           MatchManagerAction = "SET_MATCH"
	ADD_STAGED_MATCH    MatchManagerAction = "ADD_STAGED_MATCH"
	REMOVE_STAGED_MATCH MatchManagerAction = "REMOVE_STAGED_MATCH"
)

type MatchConfig struct {
	ConfigI
}

func NewMatchConfig() *MatchConfig {
	return &MatchConfig{}
}

type MatchServiceI interface {
	GetMatchById(matchId string) (*Match, error)
	GetMatchByClientKey(clientKey string) (*Match, error)
	GetStagedMatchById(matchId string) (*Match, error)
	GetStagedMatchesByClientKey(clientKey string) ([]*Match, error)

	AddMatch(match *Match) error
	StageMatchFromChallenge(challenge *Challenge) (*Match, error)
	UnstageMatch(matchId string) error
	AddMatchFromStaged(matchId string) error
	RemoveMatch(match *Match) error
	SetMatch(newMatch *Match) error

	ExecuteMove(matchId string, move *chess.Move) error
	TerminateChallenge(challenge *Challenge) error
}

type MatchService struct {
	Service[*MatchConfig]
	__dependencies__ Marker
	LoggerService    LoggerServiceI
	AuthService      AuthenticationServiceI
	SubService       SubscriptionServiceI

	__state__                Marker
	matchByMatchId           map[string]*Match
	matchIdByClientId        map[string]string
	stagedMatchById          map[string]*Match
	stagedMatchIdsByClientId map[string]*Set[string]
	mu                       sync.Mutex
}

func NewMatchService(config *MatchConfig) *MatchService {
	matchService := &MatchService{
		matchByMatchId:           make(map[string]*Match),
		matchIdByClientId:        make(map[string]string),
		stagedMatchById:          make(map[string]*Match),
		stagedMatchIdsByClientId: make(map[string]*Set[string]),
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

func (m *MatchService) GetStagedMatchById(matchId string) (*Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	match, ok := m.stagedMatchById[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not staged", matchId)
	}
	return match, nil
}

func (m *MatchService) GetMatchByClientKey(clientKey string) (*Match, error) {
	m.mu.Lock()
	matchId, ok := m.matchIdByClientId[clientKey]
	m.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("client %s not in match", clientKey)
	}
	return m.GetMatchById(matchId)
}

func (m *MatchService) GetStagedMatchesByClientKey(clientKey string) ([]*Match, error) {
	m.mu.Lock()
	matchIds, ok := m.stagedMatchIdsByClientId[clientKey]
	m.mu.Unlock()
	if !ok {
		matchIds = EmptySet[string]()
	}
	var matches []*Match
	for _, matchId := range matchIds.Flatten() {
		match, getMatchErr := m.GetStagedMatchById(matchId)
		if getMatchErr != nil {
			return nil, getMatchErr
		}
		matches = append(matches, match)
	}
	return matches, nil
}

func (m *MatchService) canStartMatchWithClientKey(clientKey string) bool {
	match, _ := m.GetMatchByClientKey(clientKey)
	return match == nil
}

func (m *MatchService) AddMatch(match *Match) error {
	m.LoggerService.Log(ENV_MATCH_MANAGER, fmt.Sprintf("adding match %s", match.Uuid))
	m.mu.Lock()
	if _, ok := m.matchByMatchId[match.Uuid]; ok {
		return fmt.Errorf("match with id %s already exists", match.Uuid)
	}
	m.mu.Unlock()
	if !m.canStartMatchWithClientKey(match.WhiteClientKey) {
		return fmt.Errorf("client %s unavailable for match", match.WhiteClientKey)
	}
	if !m.canStartMatchWithClientKey(match.BlackClientKey) {
		return fmt.Errorf("client %s unavailable for match", match.BlackClientKey)
	}
	m.mu.Lock()
	m.matchByMatchId[match.Uuid] = match
	botKey, _ := m.AuthService.GetBotKey()
	if botKey != match.WhiteClientKey {
		m.matchIdByClientId[match.WhiteClientKey] = match.Uuid
	}
	if botKey != match.BlackClientKey {
		m.matchIdByClientId[match.BlackClientKey] = match.Uuid
	}
	m.mu.Unlock()

	matchTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
	subErr := m.SubService.SubClientTo(match.WhiteClientKey, matchTopic)
	if subErr != nil {
		m.LoggerService.LogRed(ENV_MATCH_MANAGER, fmt.Sprintf("could not subscribe client %s to match topic: %s", match.WhiteClientKey, subErr))
	}
	subErr = m.SubService.SubClientTo(match.BlackClientKey, matchTopic)
	if subErr != nil {
		m.LoggerService.LogRed(ENV_MATCH_MANAGER, fmt.Sprintf("could not subscribe client %s to match topic: %s", match.BlackClientKey, subErr))
	}

	go m.startTimer(match)

	matchUpdateMsg := &Message{
		Topic:       matchTopic,
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: match,
		},
	}
	m.userClientsManager.BroadcastMessage(matchUpdateMsg)
	return nil
}

func (m *MatchService) setStagedMatch(match *Match) {
	m.mu.Lock()
	m.stagedMatchById[match.Uuid] = match
	if _, ok := m.stagedMatchIdsByClientId[match.WhiteClientKey]; !ok {
		m.stagedMatchIdsByClientId[match.WhiteClientKey] = EmptySet[string]()
	}
	if _, ok := m.stagedMatchIdsByClientId[match.BlackClientKey]; !ok {
		m.stagedMatchIdsByClientId[match.BlackClientKey] = EmptySet[string]()
	}
	whiteStagedMatchIds := m.stagedMatchIdsByClientId[match.WhiteClientKey]
	whiteStagedMatchIds.Add(match.Uuid)
	blackStagedMatchIds := m.stagedMatchIdsByClientId[match.BlackClientKey]
	blackStagedMatchIds.Add(match.Uuid)
	m.mu.Unlock()
}

func (m *MatchService) StageMatchFromChallenge(challenge *Challenge) (*Match, error) {
	if !m.canStartMatchWithClientKey(challenge.ChallengerKey) {
		return nil, fmt.Errorf("challenger client %s unavailable for match", challenge.ChallengerKey)
	}
	if !m.canStartMatchWithClientKey(challenge.ChallengedKey) {
		return nil, fmt.Errorf("challenged client %s unavailable for match", challenge.ChallengedKey)
	}
	m.LoggerService.Log(ENV_MATCH_MANAGER, fmt.Sprintf("staging match for challenger %s challenging %s", challenge.ChallengerKey, challenge.ChallengedKey))

	matchBuilder := NewMatchBuilder()
	if challenge.IsChallengerWhite {
		matchBuilder.WithWhiteClientKey(challenge.ChallengerKey)
		matchBuilder.WithBlackClientKey(challenge.ChallengedKey)
	} else if challenge.IsChallengerBlack {
		matchBuilder.WithWhiteClientKey(challenge.ChallengedKey)
		matchBuilder.WithBlackClientKey(challenge.ChallengerKey)
	} else { // challenge does not specify player colors, randomize
		matchBuilder.WithClientKeys(challenge.ChallengerKey, challenge.ChallengedKey)
	}
	matchBuilder.WithTimeControl(challenge.TimeControl)
	matchBuilder.WithTimeRemainingSec(float64(challenge.TimeControl.InitialTimeSec))
	match := matchBuilder.Build()
	m.setStagedMatch(match)
	return match, nil
}

func (m *MatchService) UnstageMatch(matchId string) error {
	m.LoggerService.Log(ENV_MATCH_MANAGER, fmt.Sprintf("unstaging match %s", matchId))
	match, getMatchErr := m.GetStagedMatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	m.mu.Lock()
	delete(m.stagedMatchById, matchId)
	if whiteStagedMatchIds, ok := m.stagedMatchIdsByClientId[match.WhiteClientKey]; ok {
		whiteStagedMatchIds.Remove(matchId)
	}
	if blackStagedMatchIds, ok := m.stagedMatchIdsByClientId[match.BlackClientKey]; ok {
		blackStagedMatchIds.Remove(matchId)
	}
	m.mu.Unlock()
	return nil
}

func (m *MatchService) AddMatchFromStaged(matchId string) error {
	stagedMatch, fetchStagedMatchErr := m.GetStagedMatchById(matchId)
	if fetchStagedMatchErr != nil {
		return fetchStagedMatchErr
	}

	whiteStagedMatches, getBlackStagedMatchesErr := m.GetStagedMatchesByClientKey(stagedMatch.WhiteClientKey)
	if getBlackStagedMatchesErr != nil {
		return fmt.Errorf("could not get staged matches for client %s: %s", stagedMatch.WhiteClientKey, getBlackStagedMatchesErr)
	}
	for _, whiteStagedMatch := range whiteStagedMatches {
		unstageErr := m.UnstageMatch(whiteStagedMatch.Uuid)
		if unstageErr != nil {
			return fmt.Errorf("could not unstage match %s: %s", whiteStagedMatch.Uuid, unstageErr.Error())
		}
	}

	blackStagedMatches, getBlackStagedMatchesErr := m.GetStagedMatchesByClientKey(stagedMatch.BlackClientKey)
	if getBlackStagedMatchesErr != nil {
		return fmt.Errorf("could not get staged matches for client %s: %s", stagedMatch.BlackClientKey, getBlackStagedMatchesErr.Error())
	}
	for _, blackStagedMatch := range blackStagedMatches {
		unstageErr := m.UnstageMatch(blackStagedMatch.Uuid)
		if unstageErr != nil {
			return fmt.Errorf("could not unstage match %s: %s", blackStagedMatch.Uuid, unstageErr.Error())
		}
	}

	addMatchErr := m.AddMatch(stagedMatch)
	if addMatchErr != nil {
		return fmt.Errorf("could not add match from staged: %s", addMatchErr.Error())
	}
	return nil
}

func (m *MatchService) RemoveMatch(match *Match) error {
	m.LoggerService.Log(ENV_MATCH_MANAGER, fmt.Sprintf("removing match %s", match.Uuid))
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.matchByMatchId[match.Uuid]; !ok {
		return fmt.Errorf("match with id %s doesn't exist", match.Uuid)
	}
	if match.WhiteClientKey != "" {
		delete(m.matchIdByClientId, match.WhiteClientKey)
	}
	if match.BlackClientKey != "" {
		delete(m.matchIdByClientId, match.BlackClientKey)
	}
	delete(m.matchByMatchId, match.Uuid)
	return nil
}

func (m *MatchService) SetMatch(newMatch *Match) error {
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

	matchUpdateMsg := &Message{
		Topic:       MessageTopic(fmt.Sprintf("match-%s", newMatch.Uuid)),
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: newMatch,
		},
	}
	m.userClientsManager.BroadcastMessage(matchUpdateMsg)
	return nil
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

	setMatchErr := m.SetMatch(newMatch)
	if setMatchErr != nil {
		return setMatchErr
	}

	if newMatch.Board.IsTerminal {
		return m.RemoveMatch(newMatch)
	} else {
		go m.startTimer(newMatch)
	}

	return nil
}

func (m *MatchService) TerminateChallenge(challenge *Challenge) error {
	m.LoggerService.Log(ENV_MATCH_MANAGER, fmt.Sprintf("failing challenge for client %s", challenge.ChallengerKey))
	stagedMatch, fetchStagedMatchErr := m.GetStagedMatchById(challenge.ChallengerKey)
	if fetchStagedMatchErr != nil {
		return fetchStagedMatchErr
	}
	_ = m.UnstageMatch(stagedMatch.Uuid)
	return nil
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
		_ = m.SetMatch(newMatch)
	}
}
