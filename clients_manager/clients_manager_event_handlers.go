package clients_manager

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

var OnCredsChanged = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with creds changed: "
	payload := event.Payload().(*auth.CredsChangedPayload)

	var keyChanged = payload.OldCreds == nil ||
		payload.OldCreds.ClientKey != payload.NewCreds.ClientKey ||
		payload.OldCreds.PrivateKey != payload.NewCreds.PrivateKey

	if keyChanged {
		sendDeps := NewSendDirectDeps(c.DirectMessage, payload.NewCreds.ClientKey)
		sendErr := SendAuth(sendDeps, payload.NewCreds.PrivateKey)
		if sendErr != nil {
			c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr)
		}
	}

	return true
}

var OnCredsVetted = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with creds vetted: "
	payload := event.Payload().(*auth.CredsVettedPayload)

	sendDeps := NewSendDirectDeps(c.DirectMessage, payload.ClientKey)
	sendErr := SendAuth(sendDeps, payload.PriKey)
	if sendErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr)
	}

	return true
}

var OnUpgradeAuthGranted = func(self ServiceI, event EventI) bool {
	c := self.(*ClientsManager)
	baseErrMsg := "could not follow up with GRANTED upgrade auth request: "
	payload := event.Payload().(*auth.RoleSwitchedPayload)

	sendDeps := NewSendDirectDeps(c.DirectMessage, payload.ClientKey)
	sendErr := SendUpgradeAuthGranted(sendDeps, payload.Role)
	if sendErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr)
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
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, sendErr)
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
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenger)")
	}
	if challengedSubErr != nil {
		c.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenged)")
	}

	sendTopicDeps := NewSendTopicDeps(c.BroadcastMessage, challenge.Topic())
	SendChallengeUpdateToAll(sendTopicDeps, challenge)
	return true
}

var OnChallengeRevoked = func(self ServiceI, event EventI) bool {
	clientManager := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeRevokedEventPayload).Challenge
	inactiveChallenge := builders.NewChallengeBuilder().FromChallenge(challenge).WithIsActive(false).Build()

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdateToAll(sendTopicDeps, inactiveChallenge)

	_ = clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	return true
}

var OnChallengeDenied = func(self ServiceI, event EventI) bool {
	clientManager := self.(*ClientsManager)
	challenge := event.Payload().(*matcher.ChallengeDeniedEventPayload).Challenge
	inactiveChallenge := builders.NewChallengeBuilder().FromChallenge(challenge).WithIsActive(false).Build()

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdateToAll(sendTopicDeps, inactiveChallenge)

	_ = clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	_ = clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	return true
}

var OnChallengeAccepted = func(s ServiceI, event EventI) bool {
	clientManager := s.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge accepted: "
	challenge := event.Payload().(*matcher.ChallengeAcceptedEventPayload).Challenge
	inactiveChallenge := builders.NewChallengeBuilder().FromChallenge(challenge).WithIsActive(false).Build()

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdateToAll(sendTopicDeps, inactiveChallenge)

	challengerSubErr := clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	challengedSubErr := clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	if challengerSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenger)")
	}
	if challengedSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenged)")
	}

	return true
}

var OnChallengeAcceptFailed = func(s ServiceI, event EventI) bool {
	clientManager := s.(*ClientsManager)
	baseErrMsg := "could not follow up on challenge accept failed: "

	challenge := event.Payload().(*matcher.ChallengeAcceptFailedEventPayload).Challenge
	if challenge == nil {
		return true
	}
	inactiveChallenge := builders.NewChallengeBuilder().FromChallenge(challenge).WithIsActive(false).Build()

	sendTopicDeps := NewSendTopicDeps(clientManager.BroadcastMessage, challenge.Topic())
	SendChallengeUpdateToAll(sendTopicDeps, inactiveChallenge)

	challengerSubErr := clientManager.SubService.UnsubClient(challenge.ChallengerKey, challenge.Topic())
	challengedSubErr := clientManager.SubService.UnsubClient(challenge.ChallengedKey, challenge.Topic())

	if challengerSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenger)")
	}
	if challengedSubErr != nil {
		clientManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, challengerSubErr, " (challenged)")
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
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, whiteSubErr, " (white)")
	}
	if blackSubErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, baseErrMsg, blackSubErr, " (black)")
	}

	deps := NewSendTopicDeps(clientsManager.BroadcastMessage, match.Topic())
	SendMatchUpdateToAll(deps, match)

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

var OnMatchUpdated = func(self ServiceI, event EventI) bool {
	clientsManager := self.(*ClientsManager)
	match := event.Payload().(*matcher.MatchUpdatedEventPayload).Match

	deps := NewSendTopicDeps(clientsManager.BroadcastMessage, match.Topic())
	SendMatchUpdateToAll(deps, match)

	return true
}

var OnMatchEnded = func(self ServiceI, event EventI) bool {
	clientsManager := self.(*ClientsManager)
	match := event.Payload().(*matcher.MatchEndedEventPayload).Match

	whiteUnsubErr := clientsManager.SubService.UnsubClient(match.WhiteClientKey, match.Topic())
	blackUnsubErr := clientsManager.SubService.UnsubClient(match.BlackClientKey, match.Topic())
	if whiteUnsubErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, "could not unsub white client from match topic", whiteUnsubErr)
	}
	if blackUnsubErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, "could not unsub black client from match topic", blackUnsubErr)
	}

	return true
}

var OnMoveFailed = func(self ServiceI, event EventI) bool {
	clientsManager := self.(*ClientsManager)
	payload := event.Payload().(*matcher.MoveFailureEventPayload)

	clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, "move failed", payload.Reason)

	sendDeps := NewSendDirectDeps(clientsManager.DirectMessage, payload.OriginClientKey)
	if sendErr := SendMoveFailed(sendDeps, payload.Move, payload.Reason); sendErr != nil {
		clientsManager.Logger.LogRed(models.ENV_CLIENT_MNGR, "could not send move failed message", sendErr)
	}

	return true
}
