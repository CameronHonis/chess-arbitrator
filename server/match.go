package server

import (
	"github.com/CameronHonis/chess"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Match struct {
	Uuid                  string       `json:"uuid"`
	Board                 *chess.Board `json:"board"`
	WhiteClientKey        string       `json:"whiteClientKey"`
	WhiteTimeRemainingSec float64      `json:"whiteTimeRemainingSec"`
	BlackClientKey        string       `json:"blackClientKey"`
	BlackTimeRemainingSec float64      `json:"blackTimeRemainingSec"`
	TimeControl           *TimeControl `json:"timeControl"`
	LastMoveTime          *time.Time   `json:"-"`
}

func NewMatch(whiteClientKey string, blackClientKey string, timeControl *TimeControl) *Match {
	matchId := uuid.New().String()
	now := time.Now()
	return &Match{
		Uuid:                  matchId,
		Board:                 chess.GetInitBoard(),
		WhiteClientKey:        whiteClientKey,
		WhiteTimeRemainingSec: float64(timeControl.InitialTimeSec),
		BlackClientKey:        blackClientKey,
		BlackTimeRemainingSec: float64(timeControl.InitialTimeSec),
		TimeControl:           timeControl,
		LastMoveTime:          &now,
	}
}

type MatchBuilder struct {
	match *Match
}

func NewMatchBuilder() *MatchBuilder {
	now := time.Now()
	return &MatchBuilder{
		match: &Match{
			Uuid:         uuid.New().String(),
			Board:        chess.GetInitBoard(),
			LastMoveTime: &now,
		},
	}
}

func (mb *MatchBuilder) WithUuid(uuid string) *MatchBuilder {
	mb.match.Uuid = uuid
	return mb
}

func (mb *MatchBuilder) WithBoard(board *chess.Board) *MatchBuilder {
	mb.match.Board = board
	return mb
}

func (mb *MatchBuilder) WithWhiteClientKey(clientKey string) *MatchBuilder {
	mb.match.WhiteClientKey = clientKey
	return mb
}

func (mb *MatchBuilder) WithWhiteTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.WhiteTimeRemainingSec = timeRemainingSec
	return mb
}

func (mb *MatchBuilder) WithBlackClientKey(clientKey string) *MatchBuilder {
	mb.match.BlackClientKey = clientKey
	return mb
}

func (mb *MatchBuilder) WithBlackTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.BlackTimeRemainingSec = timeRemainingSec
	return mb
}

func (mb *MatchBuilder) WithTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.WhiteTimeRemainingSec = timeRemainingSec
	mb.match.BlackTimeRemainingSec = timeRemainingSec
	return mb
}

func (mb *MatchBuilder) WithTimeControl(timeControl *TimeControl) *MatchBuilder {
	mb.match.TimeControl = timeControl
	return mb
}

func (mb *MatchBuilder) WithLastMoveTime(lastMoveTime *time.Time) *MatchBuilder {
	mb.match.LastMoveTime = lastMoveTime
	return mb
}

func (mb *MatchBuilder) WithClientKeys(clientAKey string, clientBKey string) *MatchBuilder {
	rand.Seed(time.Now().UnixNano())
	clientAIsWhite := rand.Intn(2) == 0
	var whiteClientKey, blackClientKey string
	if clientAIsWhite {
		whiteClientKey = clientAKey
		blackClientKey = clientBKey
	} else {
		whiteClientKey = clientBKey
		blackClientKey = clientAKey
	}
	mb.match.WhiteClientKey = whiteClientKey
	mb.match.BlackClientKey = blackClientKey
	return mb
}

func (mb *MatchBuilder) FromChallenge(challenge *Challenge) *MatchBuilder {
	mb.match = NewMatch(challenge.ChallengerKey, challenge.ChallengedKey, challenge.TimeControl)
	if challenge.IsChallengerWhite {
		mb.WithWhiteClientKey(challenge.ChallengerKey)
		mb.WithBlackClientKey(challenge.ChallengedKey)
	} else if challenge.IsChallengerBlack {
		mb.WithWhiteClientKey(challenge.ChallengedKey)
		mb.WithBlackClientKey(challenge.ChallengerKey)
	} else {
		mb.WithClientKeys(challenge.ChallengerKey, challenge.ChallengedKey)
	}
	return mb
}

func (mb *MatchBuilder) FromMatch(match *Match) *MatchBuilder {
	matchCopy := *match
	mb.match = &matchCopy
	return mb
}

func (mb *MatchBuilder) Build() *Match {
	return mb.match
}
