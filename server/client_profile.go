package server

type ClientProfile struct {
	ClientKey  Key
	Elo        int
	WinStreak  int
	LossStreak int
}

func NewClientProfile(clientKey Key, elo int) *ClientProfile {
	return &ClientProfile{
		ClientKey:  clientKey,
		Elo:        elo,
		WinStreak:  0,
		LossStreak: 0,
	}
}
