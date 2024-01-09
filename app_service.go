package main

import (
	"github.com/CameronHonis/chess-arbitrator/router_service"
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
	Service

	__dependencies__ Marker
	RouterService    *router_service.RouterService

	__state__ Marker
}

func NewAppService(config *AppConfig) *AppService {
	app := &AppService{}
	app.Service = *NewService(app, config)
	return app
}
