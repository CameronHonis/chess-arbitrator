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
}

func NewMatch(clientAKey string, clientBKey string, timeControl *TimeControl) *Match {
	rand.Seed(time.Now().UnixNano())
	clientAIsWhite := rand.Intn(2) == 0
	matchId := uuid.New().String()
	if clientAIsWhite {
		return &Match{
			Uuid:               matchId,
			Board:              chess.GetInitBoard(),
			WhiteClientId:      clientAKey,
			WhiteTimeRemaining: float64(timeControl.InitialTimeSeconds),
			BlackClientId:      clientBKey,
			BlackTimeRemaining: float64(timeControl.InitialTimeSeconds),
			TimeControl:        timeControl,
		}
	} else {
		return &Match{
			Uuid:               matchId,
			Board:              chess.GetInitBoard(),
			WhiteClientId:      clientBKey,
			WhiteTimeRemaining: float64(timeControl.InitialTimeSeconds),
			BlackClientId:      clientAKey,
			BlackTimeRemaining: float64(timeControl.InitialTimeSeconds),
			TimeControl:        timeControl,
		}
	}
}
