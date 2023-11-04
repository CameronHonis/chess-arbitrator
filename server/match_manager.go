package server

import (
	"fmt"
)

var matchManager *MatchManager

type MatchManager struct {
	matchByMatchId    map[string]*Match
	matchIdByClientId map[string]string
	stagedMatchById   map[string]*Match //only for bot matches currently
}

func GetMatchManager() *MatchManager {
	if matchManager == nil {
		matchManager = &MatchManager{
			matchByMatchId:    make(map[string]*Match),
			matchIdByClientId: make(map[string]string),
			stagedMatchById:   make(map[string]*Match),
		}
	}
	return matchManager
}

func (mm *MatchManager) AddMatch(match *Match) error {
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
	return nil
}

func (mm *MatchManager) StageMatch(match *Match) {
	mm.stagedMatchById[match.Uuid] = match
}

func (mm *MatchManager) GetStagedMatchById(matchId string) (*Match, error) {
	match, ok := mm.stagedMatchById[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not staged", matchId)
	}
	return match, nil
}

func (mm *MatchManager) UnstageMatch(matchId string) {
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
	return nil
}

func (mm *MatchManager) GetMatchById(matchId string) (*Match, error) {
	match, ok := mm.matchByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("match with id %s not found", matchId)
	}
	return match, nil
}

func (mm *MatchManager) GetMatchByClientKey(clientKey string) (*Match, error) {
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
