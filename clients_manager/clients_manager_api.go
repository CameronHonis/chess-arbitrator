package clients_manager

import "github.com/CameronHonis/chess-arbitrator/models"

type ClientWriter func(client *models.Client, msg *models.Message) error
type ErrLogger func(env string, msgs ...interface{})

type SendMessageDeps struct {
	writer ClientWriter
	client *models.Client
}

func NewSendMessageDeps(writer ClientWriter, client *models.Client) *SendMessageDeps {
	return &SendMessageDeps{writer, client}
}

func SendAuth(deps *SendMessageDeps) error {
	return deps.writer(deps.client, &models.Message{
		ContentType: models.CONTENT_TYPE_AUTH,
		Content: &models.AuthMessageContent{
			PublicKey:  deps.client.PublicKey(),
			PrivateKey: deps.client.PrivateKey(),
		},
	})
}

func SendUpgradeAuthGranted(deps *SendMessageDeps, role models.RoleName) error {
	return deps.writer(deps.client, &models.Message{
		ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
		Content: &models.UpgradeAuthGrantedMessageContent{
			UpgradedToRole: role,
		},
	})
}

func SendChallengeRequestFailed(deps *SendMessageDeps, challenge *models.Challenge, reason string) error {
	return deps.writer(deps.client, &models.Message{
		ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST_FAILED,
		Content: &models.ChallengeRequestFailedMessageContent{
			Challenge: challenge,
			Reason:    reason,
		},
	})
}

func SendChallengeRequest(deps *SendMessageDeps, challenge *models.Challenge) error {
	return deps.writer(deps.client, &models.Message{
		ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
		Content: &models.ChallengeRequestMessageContent{
			Challenge: challenge,
		},
	})
}
