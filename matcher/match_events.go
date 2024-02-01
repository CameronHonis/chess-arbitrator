package matcher

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

const (
	MATCH_CREATED         service.EventVariant = "MATCH_CREATED"
	MATCH_ENDED                                = "MATCH_ENDED"
	MATCH_UPDATED                              = "MATCH_UPDATED"
	MATCH_CREATION_FAILED                      = "MATCH_CREATION_FAILED"
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
