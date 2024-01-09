package main

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/router_service"
	"github.com/CameronHonis/chess-arbitrator/user_clients_service"
	. "github.com/CameronHonis/log"
)

func GetAppConfig() *AppConfig {
	return NewAppConfig()
}

func GetLoggerConfig() *LoggerConfig {
	logConfigBuilder := NewConfigBuilder()
	logConfigBuilder.WithDecorator(models.ENV_SERVER, WrapGreen)
	logConfigBuilder.WithDecorator(models.ENV_CLIENT, WrapBlue)
	logConfigBuilder.WithDecorator(models.ENV_MATCHMAKING, WrapMagenta)
	logConfigBuilder.WithDecorator(models.ENV_MATCH_SERVICE, WrapBrown)
	logConfigBuilder.WithDecorator(models.ENV_TIMER, WrapOrange)
	//logConfigBuilder.WithMutedEnv("server")
	//logConfigBuilder.WithMutedEnv("client")
	//logConfigBuilder.WithMutedEnv("matchmaking")
	//logConfigBuilder.WithMutedEnv("match_manager")
	//logConfigBuilder.WithMutedEnv("timer")

	return logConfigBuilder.Build()
}

func GetRouterConfig() *router_service.RouterConfig {
	return router_service.NewRouterConfig()
}

func GetUserClientsConfig() *user_clients_service.UserClientsConfig {
	return user_clients_service.NewUserClientsConfig()
}
