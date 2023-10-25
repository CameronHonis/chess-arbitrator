package server

import (
	"fmt"
)

var matchManager *MatchManager

type MatchManager struct {
	matchByMatchId    map[string]*Match
	matchIdByClientId map[string]string
}

func GetMatchManager() *MatchManager {
	if matchManager == nil {
		matchManager = &MatchManager{
			matchByMatchId:    make(map[string]*Match),
			matchIdByClientId: make(map[string]string),
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
	mm.matchIdByClientId[match.WhiteClientId] = match.Uuid
	mm.matchIdByClientId[match.BlackClientId] = match.Uuid
	matchTopic := MessageTopic(match.Uuid)
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
	HandleMessage(msg, "")
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
