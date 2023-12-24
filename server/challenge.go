package server

type Challenge struct {
	ChallengerKey     string       `json:"challengerKey"`
	ChallengedKey     string       `json:"challengedKey"`
	IsChallengerWhite bool         `json:"isChallengerWhite"`
	IsChallengerBlack bool         `json:"isChallengerBlack"`
	TimeControl       *TimeControl `json:"timeControl"`
}

func NewChallenge(challengerKey string, challengedKey string, timeControl *TimeControl) *Challenge {
	return &Challenge{
		ChallengerKey:     challengerKey,
		ChallengedKey:     challengedKey,
		IsChallengerWhite: true,
		IsChallengerBlack: false,
		TimeControl:       timeControl,
	}
}
