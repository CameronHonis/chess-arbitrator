package matcher

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

const (
	MATCH_CREATED         = "MATCH_CREATED"
	MATCH_ENDED           = "MATCH_ENDED"
	MATCH_UPDATED         = "MATCH_UPDATED"
	MATCH_CREATION_FAILED = "MATCH_CREATION_FAILED"
)

type MatchCreatedEventPayload struct {
	Match *models.Match
}

type MatchCreatedEvent struct{ Event }

func NewMatchCreatedEvent(match *models.Match) *MatchCreatedEvent {
	return &MatchCreatedEvent{
		Event: *NewEvent(MATCH_CREATED, &MatchCreatedEventPayload{
			Match: match,
		}),
	}
}

type MatchEndedEventPayload struct {
	Match *models.Match
}

type MatchEndedEvent struct{ Event }

func NewMatchEndedEvent(match *models.Match) *MatchEndedEvent {
	return &MatchEndedEvent{
		Event: *NewEvent(MATCH_ENDED, &MatchEndedEventPayload{
			Match: match,
		}),
	}
}

type MatchUpdatedEventPayload struct {
	Match *models.Match
}

type MatchUpdatedEvent struct{ Event }

func NewMatchUpdated(match *models.Match) *MatchUpdatedEvent {
	return &MatchUpdatedEvent{
		Event: *NewEvent(MATCH_UPDATED, &MatchUpdatedEventPayload{
			Match: match,
		}),
	}
}

type MatchCreationFailedEventPayload struct {
	ChallengerClientKey models.Key
	Reason              string
}

type MatchCreationFailedEvent struct{ Event }

func NewMatchCreationFailedEvent(challengerKey models.Key, reason string) *MatchCreationFailedEvent {
	return &MatchCreationFailedEvent{
		Event: *NewEvent(MATCH_CREATION_FAILED, &MatchCreationFailedEventPayload{
			ChallengerClientKey: challengerKey,
			Reason:              reason,
		}),
	}
}
