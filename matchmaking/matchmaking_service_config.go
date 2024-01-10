package matchmaking

import (
	"github.com/CameronHonis/service"
)

type MatchmakingConfig struct {
	service.ConfigI
}

func NewMatchmakingConfig() *MatchmakingConfig {
	return &MatchmakingConfig{}
}
