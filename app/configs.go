package app

import (
	cm "github.com/CameronHonis/chess-arbitrator/clients_manager"
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
	logConfigBuilder.WithDecorator(models.ENV_MATCHER_SERVICE, WrapMagenta)
	logConfigBuilder.WithDecorator(models.ENV_TIMER, WrapOrange)
	logConfigBuilder.WithDecorator(models.SUB_SERVICE, WrapOrange)
	logConfigBuilder.WithDecoratorRule(helpers.PrettyClientDecoratorRule(helpers.IsClientKey))
	//logConfigBuilder.WithMutedEnv("server")
	//logConfigBuilder.WithMutedEnv("client")
	//logConfigBuilder.WithMutedEnv("matchmaking")
	//logConfigBuilder.WithMutedEnv("match_manager")
	logConfigBuilder.WithMutedEnv(models.ENV_TIMER)

	return logConfigBuilder.Build()
}

func GetRouterConfig() *router_service.RouterServiceConfig {
	return router_service.NewRouterServiceConfig()
}

func GetClientsManagerConfig() *cm.ClientsManagerConfig {
	configBuilder := cm.NewClientsManagerConfigBuilder()
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_ECHO, cm.HandleEchoMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_JOIN_MATCHMAKING, cm.HandleFindMatchMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_SUBSCRIBE_REQUEST, cm.HandleSubscribeRequestMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST, cm.HandleRequestUpgradeAuthMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_MOVE, cm.HandleMoveMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_RESIGN_MATCH, cm.HandleResignMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_CHALLENGE_REQUEST, cm.HandleChallengePlayerMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_ACCEPT_CHALLENGE, cm.HandleAcceptChallengeMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_DECLINE_CHALLENGE, cm.HandleDeclineChallengeMessage)
	configBuilder.WithMessageHandler(models.CONTENT_TYPE_REVOKE_CHALLENGE, cm.HandleRevokeChallengeMessage)
	return configBuilder.Build()
}
