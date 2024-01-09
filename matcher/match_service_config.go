package matcher

import (
	. "github.com/CameronHonis/service"
)

type MatcherServiceConfig struct {
	ConfigI
}

func NewMatchServiceConfig() *MatcherServiceConfig {
	return &MatcherServiceConfig{}
}
