package app_service

import (
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/message_service"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	"github.com/CameronHonis/chess-arbitrator/user_clients_service"
	. "github.com/CameronHonis/log"
)

var appService *AppService

func BuildServices() *AppService {
	// init services
	appService = NewAppService(GetAppConfig())
	loggerService := NewLoggerService(GetLoggerConfig())
	routerService := router_service.NewRouterService(GetRouterConfig())
	userClientsService := user_clients_service.NewUserClientsService(GetUserClientsConfig())
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	msgService := message_service.NewMessageHandlerService(message_service.NewMessageHandlerConfig())
	subService := subscription_service.NewSubscriptionService(subscription_service.NewSubscriptionConfig())
	authService := auth_service.NewAuthenticationService(auth_service.NewAuthenticationConfig())
	matchService := matcher.NewMatcherService(matcher.NewMatchServiceConfig())
	matchmakingService := matchmaking.NewMatchmakingService(matchmaking.NewMatchmakingConfig())

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
