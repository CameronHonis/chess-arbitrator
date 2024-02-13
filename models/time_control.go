package models

import "strconv"

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

func (tc *TimeControl) Hash() string {
	return strconv.FormatInt(tc.InitialTimeSec, 10) +
		strconv.FormatInt(tc.IncrementSec, 10) +
		strconv.FormatInt(tc.TimeAfterMovesCount, 10) +
		strconv.FormatInt(tc.SecAfterMoves, 10)
}
