package clients_manager

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

const (
	CLIENT_CREATED service.EventVariant = "CLIENT_CREATED"
	CLIENT_REMOVED                      = "CLIENT_REMOVED"
	MOVE_FAILURE                        = "MOVE_FAILURE"
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

type MoveFailureEventPayload struct {
	MatchId string
	Move    *chess.Move
	Reason  string
}

type MoveFailureEvent struct{ service.Event }

func NewMoveFailureEvent(matchId string, move *chess.Move, reason string) *MoveFailureEvent {
	return &MoveFailureEvent{
		Event: *service.NewEvent(MOVE_FAILURE, &MoveFailureEventPayload{
			MatchId: matchId,
			Move:    move,
			Reason:  reason,
		}),
	}
}
