package mocks

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/set"
	. "github.com/CameronHonis/stub"
)

type MatcherServiceMock struct {
	Stubbed[matcher.MatcherService]
	ServiceMock
}

func NewMatcherServiceMock(matchService *matcher.MatcherService) *MatcherServiceMock {
	ms := &MatcherServiceMock{}
	ms.Stubbed = *NewStubbed(ms, matchService)
	ms.ServiceMock = *NewServiceMock(&matchService.Service)
	return ms
}

func (ms *MatcherServiceMock) MatchById(matchId string) (*Match, error) {
	out := ms.Call("MatchById", matchId)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatcherServiceMock) MatchByClientKey(clientKey Key) (*Match, error) {
	out := ms.Call("MatchByClientKey", clientKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatcherServiceMock) Challenges(challengerKey Key) (*Set[*Challenge], error) {
	out := ms.Call("Challenges", challengerKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Set[*Challenge]), err
}

func (ms *MatcherServiceMock) Challenge(challengerKey Key, receiverKey Key) (*Challenge, error) {
	out := ms.Call("Challenge", challengerKey, receiverKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Challenge), err
}

func (ms *MatcherServiceMock) ExecuteMove(matchId string, move *chess.Move) error {
	out := ms.Call("ExecuteMove", matchId, move)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) ChallengePlayer(challenge *Challenge) error {
	out := ms.Call("ChallengePlayer", challenge)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) AcceptChallenge(challengedKey, challengerKey Key) error {
	out := ms.Call("AcceptChallenge", challengedKey, challengerKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) RevokeChallenge(challengerKey, challengedKey Key) error {
	out := ms.Call("RevokeChallenge", challengerKey, challengedKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) DeclineChallenge(challengedKey, challengerKey Key) error {
	out := ms.Call("DeclineChallenge", challengedKey, challengerKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) AddMatch(match *Match) error {
	out := ms.Call("AddMatch", match)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) SetMatch(newMatch *Match) error {
	out := ms.Call("SetMatch", newMatch)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) RemoveMatch(match *Match) error {
	out := ms.Call("RemoveMatch", match)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (ms *MatcherServiceMock) CanStartMatchWithClientKey(clientKey Key) bool {
	out := ms.Call("CanStartMatchWithClientKey", clientKey)
	return out[0].(bool)
}

func (ms *MatcherServiceMock) StartTimer(match *Match) {
	_ = ms.Call("StartTimer", match)
}
