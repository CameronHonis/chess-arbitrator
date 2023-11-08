package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"math"
	"sync"
	"time"
)

var matchManager *MatchManager

type MatchManager struct {
	matchByMatchId        map[string]*Match
	lastMoveTimeByMatchId map[string]time.Time
	matchIdByClientId     map[string]string
	stagedMatchById       map[string]*Match //only for bot matches currently
	mu                    sync.Mutex
}

func GetMatchManager() *MatchManager {
	if matchManager == nil {
		matchManager = &MatchManager{
			matchByMatchId:        make(map[string]*Match),
			lastMoveTimeByMatchId: make(map[string]time.Time),
			matchIdByClientId:     make(map[string]string),
			stagedMatchById:       make(map[string]*Match),
		}
	}
	return matchManager
}

func (mm *MatchManager) AddMatch(match *Match) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if _, ok := mm.matchByMatchId[match.Uuid]; ok {
		return fmt.Errorf("match with id %s already exists", match.Uuid)
	}
	if _, whiteInMatch := mm.matchIdByClientId[match.WhiteClientId]; whiteInMatch {
		return fmt.Errorf("client %s (white) already in match", match.WhiteClientId)
	}
	if _, blackInMatch := mm.matchIdByClientId[match.BlackClientId]; blackInMatch {
		return fmt.Errorf("client %s (black) already in match", match.WhiteClientId)
	}
	mm.matchByMatchId[match.Uuid] = match
	mm.lastMoveTimeByMatchId[match.Uuid] = time.Now()
	if GetAuthManager().chessBotKey != match.WhiteClientId {
		mm.matchIdByClientId[match.WhiteClientId] = match.Uuid
	}
	if GetAuthManager().chessBotKey != match.BlackClientId {
		mm.matchIdByClientId[match.BlackClientId] = match.Uuid
	}
	matchTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
	subErr := GetUserClientsManager().SubscribeClientTo(match.WhiteClientId, matchTopic)
	if subErr != nil {
		GetLogManager().LogRed("matchmaking", fmt.Sprintf("could not subscribe client %s to match topic: %s", match.WhiteClientId, subErr))
	}
	subErr = GetUserClientsManager().SubscribeClientTo(match.BlackClientId, matchTopic)
	if subErr != nil {
		GetLogManager().LogRed("matchmaking", fmt.Sprintf("could not subscribe client %s to match topic: %s", match.BlackClientId, subErr))
	}
	msg := &Message{
		Topic:       matchTopic,
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: match,
		},
	}
	GetUserClientsManager().BroadcastMessage(msg)
	StartTimer(match.Uuid)
	return nil
}

func (mm *MatchManager) StageMatch(match *Match) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.stagedMatchById[match.Uuid] = match
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

func (mm *MatchManager) UnstageMatch(matchId string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	delete(mm.stagedMatchById, matchId)
}

func (mm *MatchManager) AddMatchFromStaged(matchId string) error {
	stagedMatch, fetchStagedMatchErr := mm.GetStagedMatchById(matchId)
	if fetchStagedMatchErr != nil {
		return fetchStagedMatchErr
	}
	addMatchErr := mm.AddMatch(stagedMatch)
	if addMatchErr != nil {
		return fmt.Errorf("could not add staged match with id %s: %s", matchId, addMatchErr)
	}
	mm.UnstageMatch(matchId)
	return nil
}

func (mm *MatchManager) RemoveMatch(match *Match) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if _, ok := mm.matchByMatchId[match.Uuid]; !ok {
		return fmt.Errorf("match with id %s doesn't exist", match.Uuid)
	}
	if match.WhiteClientId != "" {
		delete(mm.matchIdByClientId, match.WhiteClientId)
	}
	if match.BlackClientId != "" {
		delete(mm.matchIdByClientId, match.BlackClientId)
	}
	delete(mm.matchByMatchId, match.Uuid)
	delete(mm.lastMoveTimeByMatchId, match.Uuid)
	return nil
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

func (mm *MatchManager) GetMatchByClientKey(clientKey string) (*Match, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	matchId, ok := mm.matchIdByClientId[clientKey]
	if !ok {
		return nil, fmt.Errorf("client %s not in match", clientKey)
	}
	match, ok := mm.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not found", matchId)
	}
	return match, nil
}

func (mm *MatchManager) GetLastMoveOccurredTime(matchId string) (*time.Time, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	lastMoveTime, ok := mm.lastMoveTimeByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("last move time for match with id %s not found", matchId)
	}
	return &lastMoveTime, nil
}

func (mm *MatchManager) SetLastMoveTimeToNow(matchId string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	_, ok := mm.lastMoveTimeByMatchId[matchId]
	if !ok {
		return fmt.Errorf("last move time for match with id %s not found", matchId)
	}
	mm.lastMoveTimeByMatchId[matchId] = time.Now()
	return nil
}

func (mm *MatchManager) ExecuteMatchMove(matchId string, move *chess.Move) error {
	match, getMatchErr := mm.GetMatchById(matchId)
	if getMatchErr != nil {
		return getMatchErr
	}
	if !chess.IsLegalMove(match.Board, move) {
		return fmt.Errorf("move %v is not legal", move)
	}
	lastMoveTime, getLastMoveNanosErr := mm.GetLastMoveOccurredTime(matchId)
	if getLastMoveNanosErr != nil {
		return getLastMoveNanosErr
	}

	mm.mu.Lock()
	currTime := time.Now()
	secondsSinceLastMove := math.Max(currTime.Sub(*lastMoveTime).Seconds(), 0.1)
	if match.Board.IsWhiteTurn {
		match.WhiteTimeRemaining -= secondsSinceLastMove
	} else {
		match.BlackTimeRemaining -= secondsSinceLastMove
	}
	chess.UpdateBoardFromMove(match.Board, move)
	mm.mu.Unlock()

	matchUpdateMsg := &Message{
		Topic:       MessageTopic(fmt.Sprintf("match-%s", match.Uuid)),
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: match,
		},
	}
	GetUserClientsManager().BroadcastMessage(matchUpdateMsg)
	return nil
}

func (mm *MatchManager) CheckMatchTime(matchId string, lastMoveTime time.Time) (*time.Time, error) {
	match, getMatchErr := mm.GetMatchById(matchId)
	if getMatchErr != nil {
		return nil, getMatchErr
	}
	newLastMoveTime, getNewLastMoveTimeErr := mm.GetLastMoveOccurredTime(matchId)
	if getNewLastMoveTimeErr != nil {
		return nil, getNewLastMoveTimeErr
	}
	if lastMoveTime.Equal(*newLastMoveTime) {
		isTimeout := match.WhiteTimeRemaining == match.BlackTimeRemaining
		isTimeout = isTimeout || (match.WhiteTimeRemaining < match.BlackTimeRemaining && match.Board.IsWhiteTurn)
		isTimeout = isTimeout || (match.BlackTimeRemaining < match.WhiteTimeRemaining && !match.Board.IsWhiteTurn)
		if isTimeout {
			mm.mu.Lock()
			match.Board.IsTerminal = true
			if match.Board.IsWhiteTurn {
				match.Board.IsBlackWinner = true
				match.WhiteTimeRemaining = 0
			} else {
				match.Board.IsWhiteWinner = true
				match.BlackTimeRemaining = 0
			}
			mm.mu.Unlock()

			msg := &Message{
				Topic:       MessageTopic(fmt.Sprintf("match-%s", matchId)),
				ContentType: CONTENT_TYPE_MATCH_UPDATE,
				Content: &MatchUpdateMessageContent{
					Match: match,
				},
			}
			GetUserClientsManager().BroadcastMessage(msg)
			_ = mm.RemoveMatch(match)
		}
	}
	return newLastMoveTime, nil
}

func (mm *MatchManager) GetMatchMinTimeout(matchId string) (*time.Duration, error) {
	match, getMatchErr := mm.GetMatchById(matchId)
	if getMatchErr != nil {
		return nil, getMatchErr
	}
	minRemainingSeconds := math.Min(match.WhiteTimeRemaining, match.BlackTimeRemaining)
	minRemainingNanos := time.Duration(int64(1_000_000_000 * minRemainingSeconds))
	return &minRemainingNanos, nil
}
