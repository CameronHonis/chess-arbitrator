package matcher

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

const (
	MATCH_CREATED         service.EventVariant = "MATCH_CREATED"
	MATCH_ENDED                                = "MATCH_ENDED"
	MATCH_UPDATED                              = "MATCH_UPDATED"
	MATCH_CREATION_FAILED                      = "MATCH_CREATION_FAILED"
	MOVE_FAILURE                               = "MOVE_FAILURE"
)

type MatchCreatedEventPayload struct {
	Match *models.Match
}

type MatchCreatedEvent struct{ service.Event }

func NewMatchCreatedEvent(match *models.Match) *MatchCreatedEvent {
	return &MatchCreatedEvent{
		Event: *service.NewEvent(MATCH_CREATED, &MatchCreatedEventPayload{
			Match: match,
		}),
	}
}

type MatchEndedEventPayload struct {
	Match *models.Match
}

type MatchEndedEvent struct{ service.Event }

func NewMatchEndedEvent(match *models.Match) *MatchEndedEvent {
	return &MatchEndedEvent{
		Event: *service.NewEvent(MATCH_ENDED, &MatchEndedEventPayload{
			Match: match,
		}),
	}
}

type MatchUpdatedEventPayload struct {
	Match *models.Match
}

type MatchUpdatedEvent struct{ service.Event }

func NewMatchUpdated(match *models.Match) *MatchUpdatedEvent {
	return &MatchUpdatedEvent{
		Event: *service.NewEvent(MATCH_UPDATED, &MatchUpdatedEventPayload{
			Match: match,
		}),
	}
}

type MatchCreationFailedEventPayload struct {
	WhiteClientKey models.Key
	BlackClientKey models.Key
	Reason         string
}

type MatchCreationFailedEvent struct{ service.Event }

func NewMatchCreationFailedEvent(whiteClientKey models.Key, blackClientKey models.Key, reason string) *MatchCreationFailedEvent {
	return &MatchCreationFailedEvent{
		Event: *service.NewEvent(MATCH_CREATION_FAILED, &MatchCreationFailedEventPayload{
			WhiteClientKey: whiteClientKey,
			BlackClientKey: blackClientKey,
			Reason:         reason,
		}),
	}
}

type MoveFailureEventPayload struct {
	MatchId         string
	Move            *chess.Move
	OriginClientKey models.Key
	Reason          string
}

type MoveFailureEvent struct{ service.Event }

func NewMoveFailureEvent(matchId string, move *chess.Move, originClientKey models.Key, reason string) *MoveFailureEvent {
	return &MoveFailureEvent{
		Event: *service.NewEvent(MOVE_FAILURE, &MoveFailureEventPayload{
			MatchId:         matchId,
			Move:            move,
			OriginClientKey: originClientKey,
			Reason:          reason,
		}),
	}
}
