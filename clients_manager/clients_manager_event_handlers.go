package clients_manager

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

var OnClientCreated = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	client := event.Payload().(*ClientCreatedEventPayload).Client
	baseErrMsg := "could not send auth: "
	sendDeps := NewSendMessageDeps(c.writeMessage, client)
	sendAuthErr := SendAuth(sendDeps)
	if sendAuthErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendAuthErr.Error())
		return false
	}

	go c.listenForUserInput(client)

	return true
}

var OnUpgradeAuthGranted = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with GRANTED upgrade auth request: "
	payload := event.Payload().(*auth.AuthUpgradeGrantedPayload)
	client, clientErr := c.GetClient(payload.ClientKey)
	if clientErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, clientErr.Error())
		return false
	}

	sendDeps := NewSendMessageDeps(c.writeMessage, client)
	sendErr := SendUpgradeAuthGranted(sendDeps, payload.Role)
	if sendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr.Error())
		return false
	}
	return true
}

var OnChallengeRequestFailed = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with FAILED challenge request: "
	payload := event.Payload().(*matcher.ChallengeRequestFailedEventPayload)
	client, clientErr := c.GetClient(payload.Challenge.ChallengerKey)
	if clientErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, clientErr.Error())
		return false
	}

	sendDeps := NewSendMessageDeps(c.writeMessage, client)
	sendErr := SendChallengeRequestFailed(sendDeps, payload.Challenge, payload.Reason)
	if sendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr.Error())
		return false
	}
	return true
}

var OnChallengeCreated = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up challenge request: "
	challenge := event.Payload().(*matcher.ChallengeCreatedEventPayload).Challenge
	challengerSendErr, challengedSendErr := SendAllChallengeUpdate(c, challenge.ChallengerKey, challenge.ChallengedKey, challenge)

	if challengerSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSendErr.Error(), " (challenger)")
	}
	if challengedSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengedSendErr.Error(), " (challenged)")
	}

	return true
}

var OnChallengeRevoked = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge revoked: "
	challenge := event.Payload().(*matcher.ChallengeRevokedEventPayload).Challenge
	challengerSendErr, challengedSendErr := SendAllChallengeUpdate(c, challenge.ChallengerKey, challenge.ChallengedKey, nil)

	if challengerSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSendErr.Error(), " (challenger)")
	}
	if challengedSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengedSendErr.Error(), " (challenged)")
	}

	return true
}

var OnChallengeDenied = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge declined: "
	challenge := event.Payload().(*matcher.ChallengeDeniedEventPayload).Challenge
	challengerSendErr, challengedSendErr := SendAllChallengeUpdate(c, challenge.ChallengerKey, challenge.ChallengedKey, nil)

	if challengerSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSendErr.Error(), " (challenger)")
	}
	if challengedSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengedSendErr.Error(), " (challenged)")
	}
	return true
}

var OnChallengeAccepted = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge accepted: "
	challenge := event.Payload().(*matcher.ChallengeAcceptedEventPayload).Challenge
	challengerSendErr, challengedSendErr := SendAllChallengeUpdate(c, challenge.ChallengerKey, challenge.ChallengedKey, nil)

	if challengerSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSendErr.Error(), " (challenger)")
	}
	if challengedSendErr != nil {
		c.LogService.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengedSendErr.Error(), " (challenged)")
	}
	return true
}

var SendAllChallengeUpdate = func(c *ClientsManager, challengerKey, challengedKey models.Key, challenge *models.Challenge) (challengerSendErr error, challengedSendErr error) {
	challengerClient, challengerErr := c.GetClient(challengerKey)
	if challengerClient != nil {
		deps := NewSendMessageDeps(c.writeMessage, challengerClient)
		challengerErr = SendChallengeUpdate(deps, challenge)
	}
	challengedClient, challengedErr := c.GetClient(challengedKey)
	if challengedClient != nil {
		deps := NewSendMessageDeps(c.writeMessage, challengedClient)
		challengedErr = SendChallengeUpdate(deps, challenge)
	}
	return challengerErr, challengedErr
}
