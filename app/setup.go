package app

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	"github.com/CameronHonis/chess-arbitrator/sub_service"
	. "github.com/CameronHonis/log"
	"github.com/CameronHonis/service"
)

var appService *AppService

func BuildServices(configs ...service.ConfigI) *AppService {
	appConfig := GetAppConfig()
	loggerConfig := GetLoggerConfig()
	routerConfig := GetRouterConfig()
	clientsManagerConfig := GetClientsManagerConfig()
	subServiceConfig := sub_service.NewSubscriptionServiceConfig()
	authServiceConfig := auth.NewAuthServiceConfig()
	matchmakingServiceConfig := matchmaking.NewMatchmakingConfig()
	matcherServiceConfig := matcher.NewMatcherServiceConfig()
	for _, config := range configs {
		if _appConfig, ok := config.(*AppServiceConfig); ok {
			appConfig = _appConfig
		} else if _loggerConfig, ok := config.(*LoggerConfig); ok {
			loggerConfig = _loggerConfig
		} else if _routerConfig, ok := config.(*router_service.RouterServiceConfig); ok {
			routerConfig = _routerConfig
		} else if _clientsManagerConfig, ok := config.(*clients_manager.ClientsManagerConfig); ok {
			clientsManagerConfig = _clientsManagerConfig
		} else if _subServiceConfig, ok := config.(*sub_service.SubscriptionServiceConfig); ok {
			subServiceConfig = _subServiceConfig
		} else if _authServiceConfig, ok := config.(*auth.AuthServiceConfig); ok {
			authServiceConfig = _authServiceConfig
		} else if _matchmakingServiceConfig, ok := config.(*matchmaking.MatchmakingConfig); ok {
			matchmakingServiceConfig = _matchmakingServiceConfig
		} else if _matcherServiceConfig, ok := config.(*matcher.MatcherServiceConfig); ok {
			matcherServiceConfig = _matcherServiceConfig
		}
	}

	// init services
	appService = NewAppService(appConfig)
	loggerService := NewLoggerService(loggerConfig)
	routerService := router_service.NewRouterService(routerConfig)
	clientsManager := clients_manager.NewClientsManager(clientsManagerConfig)
	// NOTE: mixture of `get...config` and `new...config` is intentional, trying both out
	subService := sub_service.NewSubscriptionService(subServiceConfig)
	authService := auth.NewAuthenticationService(authServiceConfig)
	matchmakingService := matchmaking.NewMatchmakingService(matchmakingServiceConfig)
	matcherService := matcher.NewMatcherService(matcherServiceConfig)

	// inject dependencies
	appService.AddDependency(routerService)
	routerService.AddDependency(clientsManager)
	routerService.AddDependency(loggerService)
	clientsManager.AddDependency(loggerService)
	clientsManager.AddDependency(subService)
	clientsManager.AddDependency(authService)
	clientsManager.AddDependency(matcherService)
	clientsManager.AddDependency(matchmakingService)
	matchmakingService.AddDependency(loggerService)
	matchmakingService.AddDependency(matcherService)
	matcherService.AddDependency(loggerService)
	matcherService.AddDependency(authService)
	matcherService.AddDependency(subService)
	subService.AddDependency(authService)
	subService.AddDependency(loggerService)

	appService.Build()

	return appService
}
