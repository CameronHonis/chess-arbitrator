package models

import (
	"fmt"
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

func (c *Challenge) Topic() MessageTopic {
	return MessageTopic(fmt.Sprintf("challenge-%s", c.Uuid))
}
