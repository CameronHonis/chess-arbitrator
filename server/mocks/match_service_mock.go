package mocks

import (
	"github.com/CameronHonis/chess"
	. "github.com/CameronHonis/chess-arbitrator/match_service"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/set"
	. "github.com/CameronHonis/stub"
)

type MatchServiceMock struct {
	Stubbed[MatchService]
	ServiceMock
}

func NewMatchServiceMock(matchService *MatchService) *MatchServiceMock {
	ms := &MatchServiceMock{}
	ms.Stubbed = *NewStubbed(ms, matchService)
	ms.ServiceMock = *NewServiceMock(&matchService.Service)
	return ms
}

func (ms *MatchServiceMock) GetMatchById(matchId string) (*Match, error) {
	out := ms.Call("GetMatchById", matchId)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatchServiceMock) GetMatchByClientKey(clientKey Key) (*Match, error) {
	out := ms.Call("GetMatchByClientKey", clientKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatchServiceMock) GetChallenges(challengerKey Key) (*Set[*Challenge], error) {
	out := ms.Call("GetChallenges", challengerKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Set[*Challenge]), err
}

func (ms *MatchServiceMock) GetChallenge(challengerKey Key, receiverKey Key) (*Challenge, error) {
	out := ms.Call("GetChallenge", challengerKey, receiverKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Challenge), err
}

func (ms *MatchServiceMock) ExecuteMove(matchId string, move *chess.Move) error {
	out := ms.Call("ExecuteMove", matchId, move)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) ChallengePlayer(challenge *Challenge) error {
	out := ms.Call("ChallengePlayer", challenge)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) AcceptChallenge(challengedKey, challengerKey Key) error {
	out := ms.Call("AcceptChallenge", challengedKey, challengerKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) RevokeChallenge(challengerKey, challengedKey Key) error {
	out := ms.Call("RevokeChallenge", challengerKey, challengedKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) DeclineChallenge(challengedKey, challengerKey Key) error {
	out := ms.Call("DeclineChallenge", challengedKey, challengerKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) AddMatch(match *Match) error {
	out := ms.Call("AddMatch", match)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) SetMatch(newMatch *Match) error {
	out := ms.Call("SetMatch", newMatch)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) RemoveMatch(match *Match) error {
	out := ms.Call("RemoveMatch", match)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatchServiceMock) CanStartMatchWithClientKey(clientKey Key) bool {
	out := ms.Call("CanStartMatchWithClientKey", clientKey)
	return out[0].(bool)
}

func (ms *MatchServiceMock) StartTimer(match *Match) {
	_ = ms.Call("StartTimer", match)
}
