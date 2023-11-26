package server

type Challenge struct {
	ChallengerKey     string       `json:"challengerKey"`
	ChallengedKey     string       `json:"challengedKey"`
	IsChallengerWhite bool         `json:"isChallengerWhite"`
	TimeControl       *TimeControl `json:"timeControl"`
}
