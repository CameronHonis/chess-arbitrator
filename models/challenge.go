package models

type Challenge struct {
	ChallengerKey     Key          `json:"challengerKey"`
	ChallengedKey     Key          `json:"challengedKey"`
	IsChallengerWhite bool         `json:"isChallengerWhite"`
	IsChallengerBlack bool         `json:"isChallengerBlack"`
	TimeControl       *TimeControl `json:"timeControl"`
	BotName           string       `json:"botName"`
}

func NewChallenge(challengerKey Key, challengedKey Key, isChallengerWhite bool,
	isColorsRandom bool, timeControl *TimeControl, botName string) *Challenge {
	if isColorsRandom {

	}
	return &Challenge{
		ChallengerKey:     challengerKey,
		ChallengedKey:     challengedKey,
		IsChallengerWhite: true,
		IsChallengerBlack: false,
		TimeControl:       timeControl,
		BotName:           botName,
	}
}
