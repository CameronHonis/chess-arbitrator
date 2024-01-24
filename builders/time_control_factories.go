package builders

import "github.com/CameronHonis/chess-arbitrator/models"

func NewRapidTimeControl() *models.TimeControl {
	return &models.TimeControl{
		InitialTimeSec:      600,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}

func NewBlitzTimeControl() *models.TimeControl {
	return &models.TimeControl{
		InitialTimeSec:      300,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}

func NewBulletTimeControl() *models.TimeControl {
	return &models.TimeControl{
		InitialTimeSec:      60,
		IncrementSec:        0,
		TimeAfterMovesCount: 0,
		SecAfterMoves:       0,
	}
}
