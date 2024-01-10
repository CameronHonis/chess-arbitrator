package router_service

import "github.com/CameronHonis/service"

type RouterServiceConfig struct {
	Port uint
}

func NewRouterServiceConfig() *RouterServiceConfig {
	return &RouterServiceConfig{
		Port: 8080,
	}
}

func (rc *RouterServiceConfig) MergeWith(other service.ConfigI) service.ConfigI {
	newConfig := *(other.(*RouterServiceConfig))
	return &newConfig
}
