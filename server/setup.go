package server

import (
	. "github.com/CameronHonis/log"
)

var appService *AppService

func BuildServices() *AppService {
	// init services
	appService = NewAppService(GetAppConfig())
	loggerService := NewLoggerService(GetLoggerConfig())
	routerService := NewRouterService()

	// inject dependencies
	appService.AddDependency(loggerService)
	appService.AddDependency()

	// establish event handlers

	return appService
}
