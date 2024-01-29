package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Challenge struct {
	Uuid              string       `json:"uuid"`
	ChallengerKey     Key          `json:"challengerKey"`
	ChallengedKey     Key          `json:"challengedKey"`
	IsChallengerWhite bool         `json:"isChallengerWhite"`
	IsChallengerBlack bool         `json:"isChallengerBlack"`
	TimeControl       *TimeControl `json:"timeControl"`
	BotName           string       `json:"botName"`
}

func NewChallenge(challengerKey Key, challengedKey Key, isChallengerWhite bool,
	isColorsRandom bool, timeControl *TimeControl, botName string) *Challenge {
	challengeId := uuid.New().String()
	if isColorsRandom {

	}
	return &Challenge{
		Uuid:              challengeId,
		ChallengerKey:     challengerKey,
		ChallengedKey:     challengedKey,
		IsChallengerWhite: true,
		IsChallengerBlack: false,
		TimeControl:       timeControl,
		BotName:           botName,
	}
}

func (c *Challenge) Topic() MessageTopic {
	return MessageTopic(fmt.Sprintf("challenge-%s", c.Uuid))
}
