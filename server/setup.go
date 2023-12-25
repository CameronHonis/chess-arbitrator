package server

import (
	. "github.com/CameronHonis/log"
)

var appService *AppService

func BuildServices() *AppService {
	// init services
	appService = NewAppService(GetAppConfig())
	loggerService := NewLoggerService(GetLoggerConfig())
	routerService := NewRouterService(GetRouterConfig())
	userClientsService := NewUserClientsService(GetUserClientsConfig())
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	msgService := NewMessageHandlerService(NewMessageHandlerConfig())
	subService := NewSubscriptionService(NewSubscriptionConfig())
	authService := NewAuthenticationService(NewAuthenticationConfig())
	matchService := NewMatchService(NewMatchServiceConfig())
	matchmakingService := NewMatchmakingService(NewMatchmakingConfig())

	// inject dependencies
	appService.AddDependency(routerService)
	routerService.AddDependency(userClientsService)
	userClientsService.AddDependency(loggerService)
	userClientsService.AddDependency(msgService)
	userClientsService.AddDependency(subService)
	userClientsService.AddDependency(authService)
	msgService.AddDependency(loggerService)
	msgService.AddDependency(authService)
	msgService.AddDependency(subService)
	msgService.AddDependency(matchService)
	msgService.AddDependency(matchmakingService)
	subService.AddDependency(authService)
	matchService.AddDependency(loggerService)
	matchService.AddDependency(authService)
	matchmakingService.AddDependency(loggerService)
	matchmakingService.AddDependency(matchService)

	// establish event handlers

	return appService
}
