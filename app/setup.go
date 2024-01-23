package app

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/msg_service"
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
	userClientsService := clients_manager.NewClientsManager(GetClientsManagerConfig())
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	msgService := msg_service.NewMessageHandlerService(msg_service.NewMessageServiceConfig())
	subService := sub_service.NewSubscriptionService(sub_service.NewSubscriptionServiceConfig())
	authService := auth.NewAuthenticationService(auth.NewAuthServiceConfig())
	matchService := matcher.NewMatcherService(matcher.NewMatcherServiceConfig())
	matchmakingService := matchmaking.NewMatchmakingService(matchmaking.NewMatchmakingConfig())

	// inject dependencies
	appService.AddDependency(routerService)
	routerService.AddDependency(userClientsService)
	routerService.AddDependency(loggerService)
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
