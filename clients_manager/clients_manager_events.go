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
	Creds *models.AuthCreds
}

type ClientCreatedEvent struct{ service.Event }

func NewClientCreatedEvent(creds *models.AuthCreds) *ClientCreatedEvent {
	return &ClientCreatedEvent{
		Event: *service.NewEvent(CLIENT_CREATED, &ClientCreatedEventPayload{
			Creds: creds,
		}),
	}
}

type ClientRemovedEventPayload struct {
	ClientKey models.Key
}

type ClientRemovedEvent struct{ service.Event }

func NewClientRemovedEvent(clientKey models.Key) *ClientRemovedEvent {
	return &ClientRemovedEvent{
		Event: *service.NewEvent(CLIENT_REMOVED, &ClientRemovedEventPayload{
			ClientKey: clientKey,
		}),
	}
}
