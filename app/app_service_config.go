package app

import "github.com/CameronHonis/service"

type AppServiceConfig struct {
	service.ConfigI
}

func NewAppServiceConfig() *AppServiceConfig {
	return &AppServiceConfig{}
}
