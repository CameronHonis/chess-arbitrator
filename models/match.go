package models

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/google/uuid"
	"time"
)

type Match struct {
	Uuid                  string       `json:"uuid"`
	Board                 *chess.Board `json:"board"`
	WhiteClientKey        Key          `json:"whiteClientKey"`
	WhiteTimeRemainingSec float64      `json:"whiteTimeRemainingSec"`
	BlackClientKey        Key          `json:"blackClientKey"`
	BlackTimeRemainingSec float64      `json:"blackTimeRemainingSec"`
	TimeControl           *TimeControl `json:"timeControl"`
	LastMoveTime          *time.Time   `json:"-"`
}

func NewMatch(whiteClientKey Key, blackClientKey Key, timeControl *TimeControl) *Match {
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

func (m *Match) Topic() MessageTopic {
	return MessageTopic(fmt.Sprintf("match-%s", m.Uuid))
}
