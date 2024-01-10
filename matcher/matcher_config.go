package matcher

import (
	. "github.com/CameronHonis/service"
)

type MatcherServiceConfig struct {
	ConfigI
}

func NewMatcherServiceConfig() *MatcherServiceConfig {
	return &MatcherServiceConfig{}
}
