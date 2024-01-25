package app

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	"github.com/CameronHonis/chess-arbitrator/sub_service"
	. "github.com/CameronHonis/log"
)

var appService *AppService

func BuildServices() *AppService {
	// init services
	appService = NewAppService(GetAppConfig())
	loggerService := NewLoggerService(GetLoggerConfig())
	routerService := router_service.NewRouterService(GetRouterConfig())
	clientsManager := clients_manager.NewClientsManager(GetClientsManagerConfig())
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	subService := sub_service.NewSubscriptionService(sub_service.NewSubscriptionServiceConfig())
	authService := auth.NewAuthenticationService(auth.NewAuthServiceConfig())
	matchmakingService := matchmaking.NewMatchmakingService(matchmaking.NewMatchmakingConfig())
	matcherService := matcher.NewMatcherService(matcher.NewMatcherServiceConfig())

	// inject dependencies
	appService.AddDependency(routerService)
	routerService.AddDependency(clientsManager)
	routerService.AddDependency(loggerService)
	clientsManager.AddDependency(loggerService)
	clientsManager.AddDependency(subService)
	clientsManager.AddDependency(authService)
	clientsManager.AddDependency(matcherService)
	clientsManager.AddDependency(matchmakingService)
	subService.AddDependency(authService)
	matchmakingService.AddDependency(loggerService)
	matchmakingService.AddDependency(matcherService)
	matcherService.AddDependency(loggerService)
	matcherService.AddDependency(authService)

	// establish event handlers

	return appService
}
