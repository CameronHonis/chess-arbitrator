package matcher

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

const (
	CHALLENGE_CREATED        = "CHALLENGE_CREATED"
	CHALLENGE_DENIED         = "CHALLENGE_DENIED"
	CHALLENGE_REVOKED        = "CHALLENGE_REVOKED"
	CHALLENGE_REQUEST_FAILED = "CHALLENGE_REQUEST_FAILED"
)

type ChallengeCreatedEventPayload struct {
	Challenge *models.Challenge
}

type ChallengeCreatedEvent struct{ Event }

func NewChallengeCreatedEvent(challenge *models.Challenge) *ChallengeCreatedEvent {
	return &ChallengeCreatedEvent{
		Event: *NewEvent(CHALLENGE_CREATED, &ChallengeCreatedEventPayload{
			Challenge: challenge,
		}),
	}
}

type ChallengeDeniedEventPayload struct {
	Challenge *models.Challenge
}

type ChallengeDeniedEvent struct{ Event }

func NewChallengeCanceledEvent(challenge *models.Challenge) *ChallengeDeniedEvent {
	return &ChallengeDeniedEvent{
		Event: *NewEvent(CHALLENGE_DENIED, &ChallengeDeniedEventPayload{
			Challenge: challenge,
		}),
	}
}

type ChallengeRevokedEventPayload struct {
	Challenge *models.Challenge
}

type ChallengeRevokedEvent struct{ Event }

func NewChallengeRevokedEvent(challenge *models.Challenge) *ChallengeRevokedEvent {
	return &ChallengeRevokedEvent{
		Event: *NewEvent(CHALLENGE_REVOKED, &ChallengeRevokedEventPayload{
			Challenge: challenge,
		}),
	}
}

type ChallengeRequestFailedEventPayload struct {
	Challenge *models.Challenge
	Reason    string
}

type ChallengeRequestFailedEvent struct{ Event }

func NewChallengeRequestFailedEvent(challenge *models.Challenge, reason string) *ChallengeRequestFailedEvent {
	return &ChallengeRequestFailedEvent{
		Event: *NewEvent(CHALLENGE_REQUEST_FAILED, &ChallengeRequestFailedEventPayload{
			Challenge: challenge,
			Reason:    reason,
		}),
	}
}
