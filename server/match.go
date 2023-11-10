package server

import (
	"github.com/CameronHonis/chess"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Match struct {
	Uuid               string       `json:"uuid"`
	Board              *chess.Board `json:"board"`
	WhiteClientId      string       `json:"whiteClientId"`
	WhiteTimeRemaining float64      `json:"whiteTimeRemaining"`
	BlackClientId      string       `json:"blackClientId"`
	BlackTimeRemaining float64      `json:"blackTimeRemaining"`
	TimeControl        *TimeControl `json:"timeControl"`
	LastMoveTime       *time.Time   `json:"-"`
}

func NewMatch(clientAKey string, clientBKey string, timeControl *TimeControl) *Match {
	rand.Seed(time.Now().UnixNano())
	clientAIsWhite := rand.Intn(2) == 0
	matchId := uuid.New().String()
	var whiteClientKey, blackClientKey string
	if clientAIsWhite {
		whiteClientKey = clientAKey
		blackClientKey = clientBKey
	} else {
		whiteClientKey = clientBKey
		blackClientKey = clientAKey
	}
	now := time.Now()
	return &Match{
		Uuid:               matchId,
		Board:              chess.GetInitBoard(),
		WhiteClientId:      whiteClientKey,
		WhiteTimeRemaining: float64(timeControl.InitialTimeSeconds),
		BlackClientId:      blackClientKey,
		BlackTimeRemaining: float64(timeControl.InitialTimeSeconds),
		TimeControl:        timeControl,
		LastMoveTime:       &now,
	}
}

type MatchBuilder struct {
	match *Match
}

func NewMatchBuilder() *MatchBuilder {
	return &MatchBuilder{
		match: &Match{},
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

func (mb *MatchBuilder) WithWhiteClientId(clientId string) *MatchBuilder {
	mb.match.WhiteClientId = clientId
	return mb
}

func (mb *MatchBuilder) WithWhiteTimeRemaining(timeRemaining float64) *MatchBuilder {
	mb.match.WhiteTimeRemaining = timeRemaining
	return mb
}

func (mb *MatchBuilder) WithBlackClientId(clientId string) *MatchBuilder {
	mb.match.BlackClientId = clientId
	return mb
}

func (mb *MatchBuilder) WithBlackTimeRemaining(timeRemaining float64) *MatchBuilder {
	mb.match.BlackTimeRemaining = timeRemaining
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

func (mb *MatchBuilder) FromMatch(match *Match) *MatchBuilder {
	matchCopy := *match
	mb.match = &matchCopy
	return mb
}

func (mb *MatchBuilder) Build() *Match {
	return mb.match
}
