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
	sendDeps := NewSendDirectDeps(c.DirectMessage, client.PublicKey())
	sendAuthErr := SendAuth(sendDeps, client)
	if sendAuthErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendAuthErr.Error())
		return false
	}

	go c.listenForUserInput(client)

	return true
}

var OnUpgradeAuthGranted = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with GRANTED upgrade auth request: "
	payload := event.Payload().(*auth.AuthUpgradeGrantedPayload)

	sendDeps := NewSendDirectDeps(c.DirectMessage, payload.ClientKey)
	sendErr := SendUpgradeAuthGranted(sendDeps, payload.Role)
	if sendErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr.Error())
		return false
	}
	return true
}

var OnChallengeRequestFailed = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with FAILED challenge request: "
	payload := event.Payload().(*matcher.ChallengeRequestFailedEventPayload)

	sendDeps := NewSendDirectDeps(c.DirectMessage, payload.Challenge.ChallengerKey)
	sendErr := SendChallengeRequestFailed(sendDeps, payload.Challenge, payload.Reason)
	if sendErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr.Error())
		return false
	}
	return true
}

var OnChallengeCreated = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeCreatedEventPayload).Challenge

	sendTopicDeps := NewSendTopicDeps(c.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, challenge)
	return true
}

var OnChallengeRevoked = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeRevokedEventPayload).Challenge

	_ = c.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = c.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	sendTopicDeps := NewSendTopicDeps(c.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, challenge)

	return true
}

var OnChallengeDenied = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeDeniedEventPayload).Challenge

	_ = c.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = c.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	sendTopicDeps := NewSendTopicDeps(c.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, challenge)

	return true
}

var OnMatchCreated = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	match := event.Payload().(*matcher.MatchCreatedEventPayload).Match

	deps := NewSendTopicDeps(c.BroadcastMessage, match.Topic())
	SendMatchUpdate(deps, match)

	return true
}

var OnMatchCreationFailed = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	payload := event.Payload().(*matcher.MatchCreationFailedEventPayload)

	whiteSendDeps := NewSendDirectDeps(c.DirectMessage, payload.WhiteClientKey)
	_ = SendMatchCreationFailed(whiteSendDeps, payload.Reason, payload.WhiteClientKey, payload.BlackClientKey)

	blackSendDeps := NewSendDirectDeps(c.DirectMessage, payload.BlackClientKey)
	_ = SendMatchCreationFailed(blackSendDeps, payload.Reason, payload.WhiteClientKey, payload.BlackClientKey)

	return true
}
