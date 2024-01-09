package models

type Challenge struct {
	ChallengerKey     Key          `json:"challengerKey"`
	ChallengedKey     Key          `json:"challengedKey"`
	IsChallengerWhite bool         `json:"isChallengerWhite"`
	IsChallengerBlack bool         `json:"isChallengerBlack"`
	TimeControl       *TimeControl `json:"timeControl"`
}

func NewChallenge(challengerKey Key, challengedKey Key, timeControl *TimeControl) *Challenge {
	return &Challenge{
		ChallengerKey:     challengerKey,
		ChallengedKey:     challengedKey,
		IsChallengerWhite: true,
		IsChallengerBlack: false,
		TimeControl:       timeControl,
	}
}
