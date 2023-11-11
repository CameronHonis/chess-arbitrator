package main

import (
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/log"
)

func main() {
	configLogger()
	server.StartWSServer()
}

func configLogger() {
	logConfigBuilder := NewLogManagerConfigBuilder()
	logConfigBuilder.WithDecorator("server", WrapGreen)
	logConfigBuilder.WithDecorator("client", WrapBlue)
	logConfigBuilder.WithDecorator("matchmaking", WrapMagenta)
	logConfigBuilder.WithDecorator("match_manager", WrapBrown)
	logConfigBuilder.WithDecorator("timer", WrapOrange)
	//logConfigBuilder.WithMutedEnv("server")
	//logConfigBuilder.WithMutedEnv("client")
	//logConfigBuilder.WithMutedEnv("matchmaking")
	//logConfigBuilder.WithMutedEnv("match_manager")
	//logConfigBuilder.WithMutedEnv("timer")
	logConfigBuilder.WithMutedEnv("client_key")

	logConfig := logConfigBuilder.Build()
	GetLogManager().InjectConfig(logConfig)
}
