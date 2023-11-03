package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Match struct {
	Uuid               string
	Board              *chess.Board
	WhiteClientId      string
	WhiteTimeRemaining float64
	BlackClientId      string
	BlackTimeRemaining float64
	TimeControl        *TimeControl
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

func (m *Match) ExecuteMove(move *chess.Move) error {
	if !chess.IsLegalMove(m.Board, move) {
		return fmt.Errorf("move %v is not legal", move)
	}
	chess.UpdateBoardFromMove(m.Board, move)
	matchUpdateMsg := &Message{
		Topic:       MessageTopic(fmt.Sprintf("match-%s", m.Uuid)),
		ContentType: CONTENT_TYPE_MATCH_UPDATE,
		Content: &MatchUpdateMessageContent{
			Match: m,
		},
	}
	GetUserClientsManager().BroadcastMessage(matchUpdateMsg)
	return nil
}
