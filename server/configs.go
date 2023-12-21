package server

import (
	. "github.com/CameronHonis/log"
)

const ENV_SERVER = "server"
const ENV_CLIENT = "client"
const ENV_MATCHMAKING = "matchmaking"
const ENV_MATCH_MANAGER = "match_manager"
const ENV_TIMER = "timer"

func GetAppConfig() *AppConfig {
	return &AppConfig{}
}

func GetLoggerConfig() *LoggerConfig {
	logConfigBuilder := NewConfigBuilder()
	logConfigBuilder.WithDecorator(ENV_SERVER, WrapGreen)
	logConfigBuilder.WithDecorator(ENV_CLIENT, WrapBlue)
	logConfigBuilder.WithDecorator(ENV_MATCHMAKING, WrapMagenta)
	logConfigBuilder.WithDecorator(ENV_MATCH_MANAGER, WrapBrown)
	logConfigBuilder.WithDecorator(ENV_TIMER, WrapOrange)
	//logConfigBuilder.WithMutedEnv("server")
	//logConfigBuilder.WithMutedEnv("client")
	//logConfigBuilder.WithMutedEnv("matchmaking")
	//logConfigBuilder.WithMutedEnv("match_manager")
	//logConfigBuilder.WithMutedEnv("timer")

	return logConfigBuilder.Build()
}
