package clients_manager

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

type MessageHandler func(*ClientsManager, *models.Message) error

type ClientsManagerConfig struct {
	service.ConfigI
	handlerByContentType map[models.ContentType]MessageHandler
}

func NewClientsManagerConfig(handlersByMsgTopic map[models.ContentType]MessageHandler) *ClientsManagerConfig {
	return &ClientsManagerConfig{
		handlerByContentType: handlersByMsgTopic,
	}
}

func (c *ClientsManagerConfig) HandlerByContentType(contentType models.ContentType) MessageHandler {
	if handler := c.handlerByContentType[contentType]; handler != nil {
		return handler
	}
	return nil
}

type ClientsManagerConfigBuilder struct {
	config *ClientsManagerConfig
}

func NewClientsManagerConfigBuilder() *ClientsManagerConfigBuilder {
	handlersByMsgTopic := make(map[models.ContentType]MessageHandler)
	return &ClientsManagerConfigBuilder{
		config: NewClientsManagerConfig(handlersByMsgTopic),
	}
}

func (b *ClientsManagerConfigBuilder) WithMessageHandler(contentType models.ContentType, handler MessageHandler) *ClientsManagerConfigBuilder {
	b.config.handlerByContentType[contentType] = handler
	return b
}

func (b *ClientsManagerConfigBuilder) Build() *ClientsManagerConfig {
	return b.config
}
