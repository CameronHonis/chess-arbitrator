package server

import (
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
)

type AppConfig struct {
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (ac *AppConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*AppConfig))
	return &newConfig
}

type AppService struct {
	Service[*AppConfig]

	__dependencies__ Marker
	RouterService    *RouterService

	__state__ Marker
}

func NewAppService(config *AppConfig) *AppService {
	app := &AppService{}
	app.Service = *NewService(app, config)
	return app
}
