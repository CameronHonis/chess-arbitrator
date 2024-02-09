package clients_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
)

func HandleEchoMessage(m *ClientsManager, msg *models.Message) error {
	_, ok := msg.Content.(*models.EchoMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to EchoMessageContent")
	}
	return m.DirectMessage(msg, msg.SenderKey)
}

func HandleFindMatchMessage(m *ClientsManager, msg *models.Message) error {
	// TODO: query for elo, winStreak, lossStreak
	return m.MatchmakingService.AddClient(&models.ClientProfile{
		ClientKey:  msg.SenderKey,
		Elo:        1000,
		WinStreak:  0,
		LossStreak: 0,
	})
}

func HandleSubscribeRequestMessage(m *ClientsManager, msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.SubscribeRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestMessageContent")
	}
	subErr := m.SubService.SubClient(msg.SenderKey, msgContent.Topic)
	return subErr
}

func HandleRequestUpgradeAuthMessage(m *ClientsManager, msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.UpgradeAuthRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
	}
	return m.AuthService.UpgradeAuth(msg.SenderKey, msgContent.Role, msgContent.Secret)
}

func HandleMoveMessage(m *ClientsManager, moveMsg *models.Message) error {
	moveMsgContent, ok := moveMsg.Content.(*models.MoveMessageContent)
	if !ok {
		return fmt.Errorf("invalid move message content")
	}
	moveErr := m.MatcherService.ExecuteMove(moveMsgContent.MatchId, moveMsgContent.Move)
	if moveErr != nil {
		go m.Dispatch(NewMoveFailureEvent(moveMsgContent.MatchId, moveMsgContent.Move, moveErr.Error()))
		return nil
	}
	return nil
}

func HandleResignMessage(m *ClientsManager, resignMsg *models.Message) error {
	resignMsgContent, ok := resignMsg.Content.(*models.ResignMessageContent)
	if !ok {
		return fmt.Errorf("invalid resign message content")
	}
	return m.MatcherService.ResignMatch(resignMsgContent.MatchId, resignMsg.SenderKey)
}

func HandleChallengePlayerMessage(m *ClientsManager, challengeMsg *models.Message) error {
	challengeMsgContent, ok := challengeMsg.Content.(*models.ChallengeRequestMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge message content")
	}
	challengeErr := m.MatcherService.RequestChallenge(challengeMsgContent.Challenge)
	if challengeErr != nil {
		return challengeErr
	}
	return nil
}

func HandleAcceptChallengeMessage(m *ClientsManager, msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.AcceptChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid accept challenge message content")
	}
	acceptChallengeErr := m.MatcherService.AcceptChallenge(msgContent.ChallengerClientKey, msg.SenderKey)
	if acceptChallengeErr != nil {
		return nil
	}
	return nil
}

func HandleDeclineChallengeMessage(m *ClientsManager, msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.DeclineChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid decline challenge message content")
	}
	declineChallengeErr := m.MatcherService.DeclineChallenge(msgContent.ChallengerClientKey, msg.SenderKey)
	if declineChallengeErr != nil {
		m.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("could not decline challenge: %s", declineChallengeErr))
	}
	return nil
}

func HandleRevokeChallengeMessage(m *ClientsManager, msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.RevokeChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid revoke challenge message content")
	}
	revokeChallengeErr := m.MatcherService.RevokeChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if revokeChallengeErr != nil {
		m.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("could not revoke challenge: %s", revokeChallengeErr))
	}
	return nil
}
