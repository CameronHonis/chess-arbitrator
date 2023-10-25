package server

type TimeControl struct {
	InitialTimeSeconds  int64
	IncrementSeconds    int64
	TimeAfterMovesCount int64
	SecondsAfterMoves   int64
}
