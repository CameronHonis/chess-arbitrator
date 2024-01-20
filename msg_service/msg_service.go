package msg_service

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/sub_service"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
)

//go:generate mockgen -destination mock/msg_service_mock.go . MessageServiceI
type MessageServiceI interface {
	service.ServiceI
	HandleMessage(msg *models.Message)
}

type MessageService struct {
	service.Service

	__dependencies__   marker.Marker
	LogService         log.LoggerServiceI
	AuthService        auth.AuthenticationServiceI
	SubService         sub_service.SubscriptionServiceI
	MatcherService     matcher.MatcherServiceI
	MatchmakingService matchmaking.MatchmakingServiceI

	__state__ marker.Marker
}

func NewMessageHandlerService(config *MessageServiceConfig) *MessageService {
	messageHandler := &MessageService{}
	messageHandler.Service = *service.NewService(messageHandler, config)

	return messageHandler
}

func (m *MessageService) HandleMessage(msg *models.Message) {
	m.LogService.Log(models.ENV_CLIENT, fmt.Sprintf("handling message %s", msg))
	var handleMsgErr error
	switch msg.ContentType {
	case models.CONTENT_TYPE_ECHO:
		handleMsgErr = m.HandleEchoMessage(msg)
	case models.CONTENT_TYPE_FIND_MATCH:
		handleMsgErr = m.HandleFindMatchMessage(msg)
	case models.CONTENT_TYPE_SUBSCRIBE_REQUEST:
		handleMsgErr = m.HandleSubscribeRequestMessage(msg)
	case models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST:
		handleMsgErr = m.HandleRequestUpgradeAuthMessage(msg)
	case models.CONTENT_TYPE_MOVE:
		handleMsgErr = m.HandleMoveMessage(msg)
	case models.CONTENT_TYPE_CHALLENGE_REQUEST:
		handleMsgErr = m.HandleChallengePlayerMessage(msg)
	case models.CONTENT_TYPE_CHALLENGE_ACCEPTED:
		handleMsgErr = m.HandleAcceptChallengeMessage(msg)
	case models.CONTENT_TYPE_CHALLENGE_DECLINED:
		handleMsgErr = m.HandleDeclineChallengeMessage(msg)
	case models.CONTENT_TYPE_CHALLENGE_REVOKED:
		handleMsgErr = m.HandleRevokeChallengeMessage(msg)
	}
	if handleMsgErr != nil {
		m.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not handle message \n\t%+v\n\t%s", *msg, handleMsgErr))
	}
}

func (m *MessageService) HandleFindMatchMessage(msg *models.Message) error {
	// TODO: query for elo, winStreak, lossStreak
	return m.MatchmakingService.AddClient(&models.ClientProfile{
		ClientKey:  msg.SenderKey,
		Elo:        1000,
		WinStreak:  0,
		LossStreak: 0,
	})
}

func (m *MessageService) HandleSubscribeRequestMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.SubscribeRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestMessageContent")
	}
	subErr := m.SubService.SubClient(msg.SenderKey, msgContent.Topic)
	return subErr
}

func (m *MessageService) HandleEchoMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.EchoMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to EchoMessageContent")
	}
	go m.Dispatch(NewEchoEvent(msgContent.Message))
	return nil
}

func (m *MessageService) HandleRequestUpgradeAuthMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.UpgradeAuthRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
	}
	return m.AuthService.UpgradeAuth(msg.SenderKey, msgContent.Role, msgContent.Secret)
}

func (m *MessageService) HandleMoveMessage(moveMsg *models.Message) error {
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

func (m *MessageService) HandleChallengePlayerMessage(challengeMsg *models.Message) error {
	challengeMsgContent, ok := challengeMsg.Content.(*models.ChallengePlayerMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge message content")
	}
	challengeErr := m.MatcherService.RequestChallenge(challengeMsgContent.Challenge)
	if challengeErr != nil {
		return challengeErr
	}
	return nil
}

func (m *MessageService) HandleAcceptChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.ChallengeAcceptedMessageContent)
	if !ok {
		return fmt.Errorf("invalid accept challenge message content")
	}
	acceptChallengeErr := m.MatcherService.AcceptChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if acceptChallengeErr != nil {
		return nil
	}
	return nil
}

func (m *MessageService) HandleDeclineChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.ChallengeDeclinedMessageContent)
	if !ok {
		return fmt.Errorf("invalid decline challenge message content")
	}
	declineChallengeErr := m.MatcherService.DeclineChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if declineChallengeErr != nil {
		m.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not decline challenge: %s", declineChallengeErr))
	}
	return nil
}

func (m *MessageService) HandleRevokeChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.ChallengeRevokedMessageContent)
	if !ok {
		return fmt.Errorf("invalid revoke challenge message content")
	}
	revokeChallengeErr := m.MatcherService.RevokeChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if revokeChallengeErr != nil {
		m.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not revoke challenge: %s", revokeChallengeErr))
	}
	return nil
}
