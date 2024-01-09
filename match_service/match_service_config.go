package match_service

import (
	. "github.com/CameronHonis/service"
)

type MatchServiceConfig struct {
	ConfigI
}

func NewMatchServiceConfig() *MatchServiceConfig {
	return &MatchServiceConfig{}
}
