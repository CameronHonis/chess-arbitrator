package app

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/message_service"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	. "github.com/CameronHonis/log"
)

var appService *AppService

func BuildServices() *AppService {
	// init services
	appService = NewAppService(GetAppConfig())
	loggerService := NewLoggerService(GetLoggerConfig())
	routerService := router_service.NewRouterService(GetRouterConfig())
	userClientsService := clients_manager.NewClientsManager(GetUserClientsConfig())
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	msgService := message_service.NewMessageHandlerService(message_service.NewMessageServiceConfig())
	subService := subscription_service.NewSubscriptionService(subscription_service.NewSubscriptionServiceConfig())
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
