package server

import "github.com/CameronHonis/chess-arbitrator/models"

type ClientProfile struct {
	ClientKey  models.Key
	Elo        int
	WinStreak  int
	LossStreak int
}

func NewClientProfile(clientKey models.Key, elo int) *ClientProfile {
	return &ClientProfile{
		ClientKey:  clientKey,
		Elo:        elo,
		WinStreak:  0,
		LossStreak: 0,
	}
}
