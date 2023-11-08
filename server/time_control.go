package server

type TimeControl struct {
	InitialTimeSeconds  int64 `json:"initialTimeSeconds"`
	IncrementSeconds    int64 `json:"incrementSeconds"`
	TimeAfterMovesCount int64 `json:"timeAfterMovesCount"`
	SecondsAfterMoves   int64 `json:"secondsAfterMoves"`
}

func NewRapidTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSeconds:  600,
		IncrementSeconds:    0,
		TimeAfterMovesCount: 0,
		SecondsAfterMoves:   0,
	}
}

func NewBlitzTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSeconds:  300,
		IncrementSeconds:    0,
		TimeAfterMovesCount: 0,
		SecondsAfterMoves:   0,
	}
}

func NewBulletTimeControl() *TimeControl {
	return &TimeControl{
		InitialTimeSeconds:  60,
		IncrementSeconds:    0,
		TimeAfterMovesCount: 0,
		SecondsAfterMoves:   0,
	}
}
