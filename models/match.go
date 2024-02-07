package models

import (
	"fmt"
	"github.com/CameronHonis/chess"
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
	BotName               string       `json:"botName"`
	LastMoveTime          *time.Time   `json:"-"`
}

func (m *Match) Topic() MessageTopic {
	return MessageTopic(fmt.Sprintf("match-%s", m.Uuid))
}
