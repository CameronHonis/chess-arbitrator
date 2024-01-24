package models

type TimeControl struct {
	InitialTimeSec      int64 `json:"initialTimeSec"`
	IncrementSec        int64 `json:"incrementSec"`
	TimeAfterMovesCount int64 `json:"timeAfterMovesCount"`
	SecAfterMoves       int64 `json:"secAfterMoves"`
}

func (tc *TimeControl) Equals(other *TimeControl) bool {
	return tc.InitialTimeSec == other.InitialTimeSec &&
		tc.IncrementSec == other.IncrementSec &&
		tc.TimeAfterMovesCount == other.TimeAfterMovesCount &&
		tc.SecAfterMoves == other.SecAfterMoves
}
