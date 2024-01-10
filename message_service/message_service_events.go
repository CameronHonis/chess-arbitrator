package message_service

import "github.com/CameronHonis/chess"

import (
	"github.com/CameronHonis/service"
)

const (
	ECHO         service.EventVariant = "ECHO"
	MOVE_FAILURE                      = "MOVE_FAILURE"
)

type EchoEventPayload struct {
	Message string
}

type EchoEvent struct{ service.Event }

func NewEchoEvent(message string) *EchoEvent {
	return &EchoEvent{
		Event: *service.NewEvent(ECHO, &EchoEventPayload{
			Message: message,
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
