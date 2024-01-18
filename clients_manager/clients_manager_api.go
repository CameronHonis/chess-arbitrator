package clients_manager

import "github.com/CameronHonis/chess-arbitrator/models"

type ClientWriter func(client *models.Client, msg *models.Message) error
type ErrLogger func(env string, msgs ...interface{})

func SendAuth(writer ClientWriter, client *models.Client) error {
	return writer(client, &models.Message{
		ContentType: models.CONTENT_TYPE_AUTH,
		Content: &models.AuthMessageContent{
			PublicKey:  client.PublicKey(),
			PrivateKey: client.PrivateKey(),
		},
	})
}

func SendUpgradeAuthGranted(writer ClientWriter, client *models.Client, role models.RoleName) error {
	return writer(client, &models.Message{
		ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
		Content: &models.UpgradeAuthGrantedMessageContent{
			UpgradedToRole: role,
		},
	})
}
