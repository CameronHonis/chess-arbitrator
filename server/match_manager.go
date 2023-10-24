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
	if _, ok := mm.matchByMatchId[match.MatchId]; ok {
		return fmt.Errorf("match with id %s already exists", match.MatchId)
	}
	if _, whiteInMatch := mm.matchIdByClientId[match.WhiteClientId]; whiteInMatch {
		return fmt.Errorf("client %s (white) already in match", match.WhiteClientId)
	}
	if _, blackInMatch := mm.matchIdByClientId[match.BlackClientId]; blackInMatch {
		return fmt.Errorf("client %s (black) already in match", match.WhiteClientId)
	}
	mm.matchByMatchId[match.MatchId] = match
	mm.matchIdByClientId[match.WhiteClientId] = match.MatchId
	mm.matchIdByClientId[match.BlackClientId] = match.MatchId
	return nil
}

func (mm *MatchManager) RemoveMatch(match *Match) error {
	if _, ok := mm.matchByMatchId[match.MatchId]; !ok {
		return fmt.Errorf("match with id %s doesn't exist", match.MatchId)
	}
	if match.WhiteClientId != "" {
		delete(mm.matchIdByClientId, match.WhiteClientId)
	}
	if match.BlackClientId != "" {
		delete(mm.matchIdByClientId, match.BlackClientId)
	}
	delete(mm.matchByMatchId, match.MatchId)
	return nil
}
