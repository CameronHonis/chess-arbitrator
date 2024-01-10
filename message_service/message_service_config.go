package message_service

import (
	"github.com/CameronHonis/service"
)

type MessageServiceConfig struct {
	service.ConfigI
}

func NewMessageServiceConfig() *MessageServiceConfig {
	return &MessageServiceConfig{}
}
