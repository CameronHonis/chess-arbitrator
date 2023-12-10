package server

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

func NewRapidTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSec:      600,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}

func NewBlitzTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSec:      300,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}

func NewBulletTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSec:      60,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}
