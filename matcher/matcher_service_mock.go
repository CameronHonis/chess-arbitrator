package matcher

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/mocks"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/set"
	. "github.com/CameronHonis/stub"
)

type MatcherServiceMock struct {
	Stubbed[MatcherService]
	mocks.ServiceMock
}

func NewMatcherServiceMock(matchService *MatcherService) *MatcherServiceMock {
	ms := &MatcherServiceMock{}
	ms.Stubbed = *NewStubbed(ms, matchService)
	ms.ServiceMock = *mocks.NewServiceMock(&matchService.Service)
	return ms
}

func (ms *MatcherServiceMock) GetMatchById(matchId string) (*Match, error) {
	out := ms.Call("GetMatchById", matchId)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatcherServiceMock) GetMatchByClientKey(clientKey Key) (*Match, error) {
	out := ms.Call("GetMatchByClientKey", clientKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Match), err
}

func (ms *MatcherServiceMock) GetChallenges(challengerKey Key) (*Set[*Challenge], error) {
	out := ms.Call("GetChallenges", challengerKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(*Set[*Challenge]), err
}

func (ms *MatcherServiceMock) GetChallenge(challengerKey Key, receiverKey Key) (*Challenge, error) {
	out := ms.Call("GetChallenge", challengerKey, receiverKey)
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
