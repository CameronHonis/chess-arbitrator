package server

import (
	"github.com/CameronHonis/chess-arbitrator/chess"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Match struct {
	MatchId            string
	board              *chess.Board
	WhiteClientId      string
	WhiteTimeRemaining float64
	BlackClientId      string
	BlackTimeRemaining float64
}

func NewMatch(clientId string) *Match {
	rand.Seed(time.Now().UnixNano())
	clientIsWhite := rand.Intn(2) == 0
	matchId := uuid.New().String()
	if clientIsWhite {
		return &Match{
			MatchId:            matchId,
			board:              chess.GetInitBoard(),
			WhiteClientId:      clientId,
			WhiteTimeRemaining: 300,
			BlackClientId:      "",
			BlackTimeRemaining: 300,
		}
	} else {
		return &Match{
			MatchId:            matchId,
			board:              chess.GetInitBoard(),
			WhiteClientId:      "",
			WhiteTimeRemaining: 300,
			BlackClientId:      clientId,
			BlackTimeRemaining: 300,
		}
	}
}
