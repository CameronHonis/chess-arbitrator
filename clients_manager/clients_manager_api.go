package clients_manager

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
)

type DirectMessageFn func(msg *models.Message, clientKey models.Key) error
type ErrLogger func(env string, msgs ...interface{})

type SendDirectDeps struct {
	writer    DirectMessageFn
	clientKey models.Key
}

func NewSendDirectDeps(writer DirectMessageFn, clientKey models.Key) *SendDirectDeps {
	return &SendDirectDeps{writer, clientKey}
}

func SendAuth(deps *SendDirectDeps, priKey models.Key) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_AUTH,
		Content: &models.AuthMessageContent{
			PublicKey:  deps.clientKey,
			PrivateKey: priKey,
		},
	}, deps.clientKey)
}

func SendInvalidAuth(deps *SendDirectDeps) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_INVALID_AUTH,
		Content:     &models.NoMessageContent{},
	}, deps.clientKey)
}

func SendUpgradeAuthGranted(deps *SendDirectDeps, role models.RoleName) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
		Content: &models.UpgradeAuthGrantedMessageContent{
			UpgradedToRole: role,
		},
	}, deps.clientKey)
}

func SendChallengeRequestFailed(deps *SendDirectDeps, challenge *models.Challenge, reason string) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST_FAILED,
		Content: &models.ChallengeRequestFailedMessageContent{
			Challenge: challenge,
			Reason:    reason,
		},
	}, deps.clientKey)
}

func SendMatchCreationFailed(deps *SendDirectDeps, reason string, whiteClientKey, blackClientKey models.Key) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_MATCH_CREATION_FAILED,
		Content: &models.MatchCreationFailedMessageContent{
			WhiteClientKey: whiteClientKey,
			BlackClientKey: blackClientKey,
			Reason:         reason,
		},
	}, deps.clientKey)
}

func SendMatchUpdate(deps *SendDirectDeps, match *models.Match) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_MATCH_UPDATED,
		Content: &models.MatchUpdateMessageContent{
			Match: match,
		},
	}, deps.clientKey)
}

func SendMoveFailed(deps *SendDirectDeps, move *chess.Move, reason string) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_MOVE_FAILED,
		Content: &models.MoveFailedMessageContent{
			Move:   move,
			Reason: reason,
		},
	}, deps.clientKey)
}

func SendChallengeUpdate(deps *SendDirectDeps, challenge *models.Challenge) error {
	return deps.writer(&models.Message{
		ContentType: models.CONTENT_TYPE_CHALLENGE_UPDATED,
		Content: &models.ChallengeUpdatedMessageContent{
			Challenge: challenge,
		},
	}, deps.clientKey)
}

type BroadcastMessageFn func(msg *models.Message)

type SendTopicDeps struct {
	writer BroadcastMessageFn
	topic  models.MessageTopic
}

func NewSendTopicDeps(writer BroadcastMessageFn, topic models.MessageTopic) *SendTopicDeps {
	return &SendTopicDeps{writer, topic}
}

func SendChallengeUpdateToAll(deps *SendTopicDeps, challenge *models.Challenge) {
	deps.writer(&models.Message{
		Topic:       deps.topic,
		ContentType: models.CONTENT_TYPE_CHALLENGE_UPDATED,
		Content: &models.ChallengeUpdatedMessageContent{
			Challenge: challenge,
		},
	})
}

func SendMatchUpdateToAll(deps *SendTopicDeps, match *models.Match) {
	deps.writer(&models.Message{
		Topic:       deps.topic,
		ContentType: models.CONTENT_TYPE_MATCH_UPDATED,
		Content: &models.MatchUpdateMessageContent{
			Match: match,
		},
	})
}
