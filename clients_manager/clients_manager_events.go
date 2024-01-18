package clients_manager

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

const (
	CLIENT_CREATED service.EventVariant = "CLIENT_CREATED"
	CLIENT_REMOVED                      = "CLIENT_REMOVED"
)

type ClientCreatedEventPayload struct {
	Client *models.Client
}

type ClientCreatedEvent struct{ service.Event }

func NewClientCreatedEvent(client *models.Client) *ClientCreatedEvent {
	return &ClientCreatedEvent{
		Event: *service.NewEvent(CLIENT_CREATED, &ClientCreatedEventPayload{
			Client: client,
		}),
	}
}

type ClientRemovedEventPayload struct {
	Client *models.Client
}

type ClientRemovedEvent struct{ service.Event }

func NewClientRemovedEvent(client *models.Client) *ClientRemovedEvent {
	return &ClientRemovedEvent{
		Event: *service.NewEvent(CLIENT_REMOVED, &ClientRemovedEventPayload{
			Client: client,
		}),
	}
}
