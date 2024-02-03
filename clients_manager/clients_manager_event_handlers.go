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
	baseErrMsg := "could not follow up on challenge created: "
	challenge := event.Payload().(*matcher.ChallengeCreatedEventPayload).Challenge

	challengerSubErr := c.SubService.SubClient(challenge.ChallengerKey, challenge.Topic())
	challengedSubErr := c.SubService.SubClient(challenge.ChallengedKey, challenge.Topic())

	if challengerSubErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr.Error(), " (challenger)")
	}
	if challengedSubErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr.Error(), " (challenged)")
	}

	sendTopicDeps := NewSendTopicDeps(c.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, challenge)
	return true
}

var OnChallengeRevoked = func(self ServiceI, event EventI) bool {
	clientManager := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeRevokedEventPayload).Challenge

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, nil)

	_ = clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	return true
}

var OnChallengeDenied = func(self ServiceI, event EventI) bool {
	clientManager := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeDeniedEventPayload).Challenge

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, nil)

	_ = clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	return true
}

var OnChallengeAccepted = func(s ServiceI, event EventI) bool {
	clientManager := s.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge accepted: "
	challenge := event.Payload().(*matcher.ChallengeAcceptedEventPayload).Challenge

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdate(sendTopicDeps, nil)

	challengerSubErr := clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	challengedSubErr := clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	if challengerSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr.Error(), " (challenger)")
	}
	if challengedSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr.Error(), " (challenged)")
	}

	return true
}

var OnMatchCreated = func(self ServiceI, event EventI) bool {
	clientsManager := self.(*ClientsManager)
	baseErrMsg := "could not follow up on match created: "
	match := event.Payload().(*matcher.MatchCreatedEventPayload).Match

	whiteSubErr := clientsManager.SubService.SubClient(match.WhiteClientKey, match.Topic())
	blackSubErr := clientsManager.SubService.SubClient(match.BlackClientKey, match.Topic())

	if whiteSubErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, whiteSubErr.Error(), " (white)")
	}
	if blackSubErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, blackSubErr.Error(), " (black)")
	}

	deps := NewSendTopicDeps(clientsManager.BroadcastMessage, match.Topic())
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
