package server

import (
	"github.com/CameronHonis/chess"
	. "github.com/CameronHonis/log"
	"time"
)

type TimerI interface {
	Start(match *Match)
}

var timer *Timer

type Timer struct {
	logManager   LogManagerI
	matchManager MatchManagerI
}

func GetTimer() *Timer {
	if timer != nil {
		return timer
	}
	timer = &Timer{} // null service to prevent infinite recursion
	timer = &Timer{
		logManager:   GetLogManager(),
		matchManager: GetMatchManager(),
	}
	return timer
}

func (t *Timer) Start(match *Match) {
	var waitTime time.Duration
	if match.Board.IsWhiteTurn {
		waitTime = time.Duration(match.WhiteTimeRemaining) * time.Second
	} else {
		waitTime = time.Duration(match.BlackTimeRemaining) * time.Second
	}

	time.Sleep(waitTime)
	currMatch, _ := GetMatchManager().GetMatchById(match.Uuid)
	if currMatch == nil {
		GetLogManager().LogRed(ENV_TIMER, "match not found")
		return
	}
	if currMatch.LastMoveTime.Equal(*match.LastMoveTime) {
		matchBuilder := NewMatchBuilder().FromMatch(match)
		boardBuilder := chess.NewBoardBuilder().FromBoard(match.Board)
		boardBuilder.WithIsTerminal(true)
		if match.Board.IsWhiteTurn {
			matchBuilder.WithWhiteTimeRemaining(0)
			boardBuilder.WithIsBlackWinner(true)
		} else {
			matchBuilder.WithBlackTimeRemaining(0)
			boardBuilder.WithIsWhiteWinner(true)
		}
		matchBuilder.WithBoard(boardBuilder.Build())
		newMatch := matchBuilder.Build()
		_ = GetMatchManager().SetMatch(newMatch)
	}
}
