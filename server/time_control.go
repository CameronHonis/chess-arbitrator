package server

type TimeControl struct {
	InitialTimeSeconds  int64 `json:"initialTimeSeconds"`
	IncrementSeconds    int64 `json:"incrementSeconds"`
	TimeAfterMovesCount int64 `json:"timeAfterMovesCount"`
	SecondsAfterMoves   int64 `json:"secondsAfterMoves"`
}
