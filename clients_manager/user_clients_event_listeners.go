package clients_manager

import "github.com/CameronHonis/chess-arbitrator/models"

type ClientWriter func(client *models.Client, msg *models.Message) error
type ErrLogger func(env string, msgs ...interface{})

func OnClientCreated(writer ClientWriter, logErr ErrLogger, client *models.Client) {
	writeErr := writer(client, &models.Message{
		ContentType: models.CONTENT_TYPE_AUTH,
		Content: &models.AuthMessageContent{
			PublicKey:  client.PublicKey(),
			PrivateKey: client.PrivateKey(),
		},
	})
	if writeErr != nil {
		logErr(models.ENV_CLIENT_MNGR, "could not send auth: ", writeErr.Error())
	}
}
