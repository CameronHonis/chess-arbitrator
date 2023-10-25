package server

type ClientProfile struct {
	ClientKey  string
	Elo        int
	WinStreak  int
	LossStreak int
}

func NewClientProfile(clientKey string, elo int) *ClientProfile {
	return &ClientProfile{
		ClientKey:  clientKey,
		Elo:        elo,
		WinStreak:  0,
		LossStreak: 0,
	}
}
