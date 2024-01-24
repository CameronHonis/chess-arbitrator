package app

import (
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/helpers"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	. "github.com/CameronHonis/log"
)

func GetAppConfig() *AppServiceConfig {
	return NewAppServiceConfig()
}

func GetLoggerConfig() *LoggerConfig {
	logConfigBuilder := NewConfigBuilder()
	logConfigBuilder.WithDecorator(models.ENV_SERVER, WrapGreen)
	logConfigBuilder.WithDecorator(models.ENV_CLIENT, WrapBlue)
	logConfigBuilder.WithDecorator(models.ENV_MATCHMAKING, WrapMagenta)
	logConfigBuilder.WithDecorator(models.ENV_MATCHER_SERVICE, WrapBrown)
	logConfigBuilder.WithDecorator(models.ENV_TIMER, WrapOrange)
	logConfigBuilder.WithDecoratorRule(helpers.PrettyClientDecoratorRule(helpers.IsClientKey))
	//logConfigBuilder.WithMutedEnv("server")
	//logConfigBuilder.WithMutedEnv("client")
	//logConfigBuilder.WithMutedEnv("matchmaking")
	//logConfigBuilder.WithMutedEnv("match_manager")
	//logConfigBuilder.WithMutedEnv("timer")

	return logConfigBuilder.Build()
}

func GetRouterConfig() *router_service.RouterServiceConfig {
	return router_service.NewRouterServiceConfig()
}

func GetClientsManagerConfig() *clients_manager.ClientsManagerConfig {
	return clients_manager.NewClientsManagerConfig(true)
}
