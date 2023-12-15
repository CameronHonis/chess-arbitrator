package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	. "github.com/CameronHonis/log"
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

type MatchManagerI interface {
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

var matchManager *MatchManager

type MatchManager struct {
	logManager          LogManagerI
	userClientsManager  UserClientsManagerI
	authManager         AuthManagerI
	subscriptionManager SubscriptionManagerI
	timer               TimerI

	sideEffects []func()

	matchByMatchId           map[string]*Match
	matchIdByClientId        map[string]string
	stagedMatchById          map[string]*Match
	stagedMatchIdsByClientId map[string]*Set[string]
	mu                       sync.Mutex
}

func GetMatchManager() *MatchManager {
	if matchManager != nil {
		return matchManager
	}
	matchManager = &MatchManager{} // null service to prevent infinite recursion
	matchManager = &MatchManager{
		logManager:          GetLogManager(),
		userClientsManager:  GetUserClientsManager(),
		authManager:         GetAuthManager(),
		subscriptionManager: GetSubscriptionManager(),
		timer:               GetTimer(),

		matchByMatchId:           make(map[string]*Match),
		matchIdByClientId:        make(map[string]string),
		stagedMatchById:          make(map[string]*Match),
		stagedMatchIdsByClientId: make(map[string]*Set[string]),
	}
	return matchManager
}

func (mm *MatchManager) GetMatchById(matchId string) (*Match, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	match, ok := mm.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not found", matchId)
	}
	return match, nil
}

func (mm *MatchManager) GetStagedMatchById(matchId string) (*Match, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	match, ok := mm.stagedMatchById[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not staged", matchId)
	}
	return match, nil
}

func (mm *MatchManager) GetMatchByClientKey(clientKey string) (*Match, error) {
	mm.mu.Lock()
	matchId, ok := mm.matchIdByClientId[clientKey]
	mm.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("client %s not in match", clientKey)
	}
	return mm.GetMatchById(matchId)
}

func (mm *MatchManager) GetStagedMatchesByClientKey(clientKey string) ([]*Match, error) {
	mm.mu.Lock()
	matchIds, ok := mm.stagedMatchIdsByClientId[clientKey]
	mm.mu.Unlock()
	if !ok {
		matchIds = EmptySet[string]()
	}
	var matches []*Match
	for _, matchId := range matchIds.Flatten() {
		match, getMatchErr := mm.GetStagedMatchById(matchId)
		if getMatchErr != nil {
			return nil, getMatchErr
		}
		matches = append(matches, match)
	}
	return matches, nil
}

func (mm *MatchManager) canStartMatchWithClientKey(clientKey string) bool {
	match, _ := mm.GetMatchByClientKey(clientKey)
	return match == nil
}

func (mm *MatchManager) AddMatch(match *Match) error {
	mm.logManager.Log(ENV_MATCH_MANAGER, fmt.Sprintf("adding match %s", match.Uuid))
	mm.mu.Lock()
	if _, ok := mm.matchByMatchId[match.Uuid]; ok {
		return fmt.Errorf("match with id %s already exists", match.Uuid)
	}
	mm.mu.Unlock()
	if !mm.canStartMatchWithClientKey(match.WhiteClientKey) {
		return fmt.Errorf("client %s unavailable for match", match.WhiteClientKey)
	}
	if !mm.canStartMatchWithClientKey(match.BlackClientKey) {
		return fmt.Errorf("client %s unavailable for match", match.BlackClientKey)
	}
	mm.mu.Lock()
	mm.matchByMatchId[match.Uuid] = match
	botKey, _ := mm.authManager.GetBotKey()
	if botKey != match.WhiteClientKey {
		mm.matchIdByClientId[match.WhiteClientKey] = match.Uuid
	}
	if botKey != match.BlackClientKey {
		mm.matchIdByClientId[match.BlackClientKey] = match.Uuid
	}
	mm.mu.Unlock()

	matchTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
	subErr := mm.subscriptionManager.SubClientTo(match.WhiteClientKey, matchTopic)
	if subErr != nil {
		mm.logManager.LogRed(ENV_MATCH_MANAGER, fmt.Sprintf("could not subscribe client %s to match topic: %s", match.WhiteClientKey, subErr))
	}
	subErr = mm.subscriptionManager.SubClientTo(match.BlackClientKey, matchTopic)
	if subErr != nil {
		mm.logManager.LogRed(ENV_MATCH_MANAGER, fmt.Sprintf("could not subscribe client %s to match topic: %s", match.BlackClientKey, subErr))
	}

	go mm.timer.Start(match)

	matchUpdateMsg := &Message{
		Topic:       matchTopic,
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: match,
		},
	}
	mm.userClientsManager.BroadcastMessage(matchUpdateMsg)
	return nil
}

func (mm *MatchManager) setStagedMatch(match *Match) {
	mm.mu.Lock()
	mm.stagedMatchById[match.Uuid] = match
	if _, ok := mm.stagedMatchIdsByClientId[match.WhiteClientKey]; !ok {
		mm.stagedMatchIdsByClientId[match.WhiteClientKey] = EmptySet[string]()
	}
	if _, ok := mm.stagedMatchIdsByClientId[match.BlackClientKey]; !ok {
		mm.stagedMatchIdsByClientId[match.BlackClientKey] = EmptySet[string]()
	}
	whiteStagedMatchIds := mm.stagedMatchIdsByClientId[match.WhiteClientKey]
	whiteStagedMatchIds.Add(match.Uuid)
	blackStagedMatchIds := mm.stagedMatchIdsByClientId[match.BlackClientKey]
	blackStagedMatchIds.Add(match.Uuid)
	mm.mu.Unlock()
}

func (mm *MatchManager) StageMatchFromChallenge(challenge *Challenge) (*Match, error) {
	if !mm.canStartMatchWithClientKey(challenge.ChallengerKey) {
		return nil, fmt.Errorf("challenger client %s unavailable for match", challenge.ChallengerKey)
	}
	if !mm.canStartMatchWithClientKey(challenge.ChallengedKey) {
		return nil, fmt.Errorf("challenged client %s unavailable for match", challenge.ChallengedKey)
	}
	mm.logManager.Log(ENV_MATCH_MANAGER, fmt.Sprintf("staging match for challenger %s challenging %s", challenge.ChallengerKey, challenge.ChallengedKey))

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
	mm.setStagedMatch(match)
	return match, nil
}

func (mm *MatchManager) UnstageMatch(matchId string) error {
	mm.logManager.Log(ENV_MATCH_MANAGER, fmt.Sprintf("unstaging match %s", matchId))
	match, getMatchErr := mm.GetStagedMatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	mm.mu.Lock()
	delete(mm.stagedMatchById, matchId)
	if whiteStagedMatchIds, ok := mm.stagedMatchIdsByClientId[match.WhiteClientKey]; ok {
		whiteStagedMatchIds.Remove(matchId)
	}
	if blackStagedMatchIds, ok := mm.stagedMatchIdsByClientId[match.BlackClientKey]; ok {
		blackStagedMatchIds.Remove(matchId)
	}
	mm.mu.Unlock()
	return nil
}

func (mm *MatchManager) AddMatchFromStaged(matchId string) error {
	stagedMatch, fetchStagedMatchErr := mm.GetStagedMatchById(matchId)
	if fetchStagedMatchErr != nil {
		return fetchStagedMatchErr
	}

	whiteStagedMatches, getBlackStagedMatchesErr := mm.GetStagedMatchesByClientKey(stagedMatch.WhiteClientKey)
	if getBlackStagedMatchesErr != nil {
		return fmt.Errorf("could not get staged matches for client %s: %s", stagedMatch.WhiteClientKey, getBlackStagedMatchesErr)
	}
	for _, whiteStagedMatch := range whiteStagedMatches {
		unstageErr := mm.UnstageMatch(whiteStagedMatch.Uuid)
		if unstageErr != nil {
			return fmt.Errorf("could not unstage match %s: %s", whiteStagedMatch.Uuid, unstageErr.Error())
		}
	}

	blackStagedMatches, getBlackStagedMatchesErr := mm.GetStagedMatchesByClientKey(stagedMatch.BlackClientKey)
	if getBlackStagedMatchesErr != nil {
		return fmt.Errorf("could not get staged matches for client %s: %s", stagedMatch.BlackClientKey, getBlackStagedMatchesErr.Error())
	}
	for _, blackStagedMatch := range blackStagedMatches {
		unstageErr := mm.UnstageMatch(blackStagedMatch.Uuid)
		if unstageErr != nil {
			return fmt.Errorf("could not unstage match %s: %s", blackStagedMatch.Uuid, unstageErr.Error())
		}
	}

	addMatchErr := mm.AddMatch(stagedMatch)
	if addMatchErr != nil {
		return fmt.Errorf("could not add match from staged: %s", addMatchErr.Error())
	}
	return nil
}

func (mm *MatchManager) RemoveMatch(match *Match) error {
	mm.logManager.Log(ENV_MATCH_MANAGER, fmt.Sprintf("removing match %s", match.Uuid))
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if _, ok := mm.matchByMatchId[match.Uuid]; !ok {
		return fmt.Errorf("match with id %s doesn't exist", match.Uuid)
	}
	if match.WhiteClientKey != "" {
		delete(mm.matchIdByClientId, match.WhiteClientKey)
	}
	if match.BlackClientKey != "" {
		delete(mm.matchIdByClientId, match.BlackClientKey)
	}
	delete(mm.matchByMatchId, match.Uuid)
	return nil
}

func (mm *MatchManager) SetMatch(newMatch *Match) error {
	oldMatch, fetchCurrMatchErr := mm.GetMatchById(newMatch.Uuid)
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
	mm.mu.Lock()
	mm.matchByMatchId[newMatch.Uuid] = newMatch
	mm.mu.Unlock()

	matchUpdateMsg := &Message{
		Topic:       MessageTopic(fmt.Sprintf("match-%s", newMatch.Uuid)),
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: newMatch,
		},
	}
	mm.userClientsManager.BroadcastMessage(matchUpdateMsg)
	return nil
}

func (mm *MatchManager) ExecuteMove(matchId string, move *chess.Move) error {
	match, getMatchErr := mm.GetMatchById(matchId)
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

	setMatchErr := mm.SetMatch(newMatch)
	if setMatchErr != nil {
		return setMatchErr
	}

	if newMatch.Board.IsTerminal {
		return mm.RemoveMatch(newMatch)
	} else {
		go mm.timer.Start(newMatch)
	}

	return nil
}

func (mm *MatchManager) TerminateChallenge(challenge *Challenge) error {
	mm.logManager.Log(ENV_MATCH_MANAGER, fmt.Sprintf("failing challenge for client %s", challenge.ChallengerKey))
	stagedMatch, fetchStagedMatchErr := mm.GetStagedMatchById(challenge.ChallengerKey)
	if fetchStagedMatchErr != nil {
		return fetchStagedMatchErr
	}
	mm.UnstageMatch(stagedMatch.Uuid)
	return nil
}
