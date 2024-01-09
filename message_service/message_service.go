package message_service

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
)

const (
	ECHO         EventVariant = "ECHO"
	MOVE_FAILURE              = "MOVE_FAILURE"
)

type EchoEventPayload struct {
	Message string
}

type EchoEvent struct{ Event }

func NewEchoEvent(message string) *EchoEvent {
	return &EchoEvent{
		Event: *NewEvent(ECHO, &EchoEventPayload{
			Message: message,
		}),
	}
}

type MoveFailureEventPayload struct {
	MatchId string
	Move    *chess.Move
	Reason  string
}

type MoveFailureEvent struct{ Event }

func NewMoveFailureEvent(matchId string, move *chess.Move, reason string) *MoveFailureEvent {
	return &MoveFailureEvent{
		Event: *NewEvent(MOVE_FAILURE, &MoveFailureEventPayload{
			MatchId: matchId,
			Move:    move,
			Reason:  reason,
		}),
	}
}

type MessageHandlerConfig struct {
}

func NewMessageHandlerConfig() *MessageHandlerConfig {
	return &MessageHandlerConfig{}
}

func (mhc *MessageHandlerConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*MessageHandlerConfig))
	return &newConfig
}

type MessageServiceI interface {
	ServiceI
	HandleMessage(msg *models.Message)
}

type MessageService struct {
	Service

	__dependencies__      Marker
	LoggerService         LoggerServiceI
	AuthenticationService auth_service.AuthenticationServiceI
	SubscriptionService   subscription_service.SubscriptionServiceI
	MatchService          matcher.MatchServiceI
	MatchmakingService    matchmaking.MatchmakingServiceI

	__state__ Marker
}

func NewMessageHandlerService(config *MessageHandlerConfig) *MessageService {
	messageHandler := &MessageService{}
	messageHandler.Service = *NewService(messageHandler, config)

	return messageHandler
}

func (m *MessageService) HandleMessage(msg *models.Message) {
	m.LoggerService.Log(models.ENV_CLIENT, fmt.Sprintf("handling message %s", msg))
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
	case models.CONTENT_TYPE_CHALLENGE_PLAYER:
		handleMsgErr = m.HandleChallengePlayerMessage(msg)
	case models.CONTENT_TYPE_ACCEPT_CHALLENGE:
		handleMsgErr = m.HandleAcceptChallengeMessage(msg)
	case models.CONTENT_TYPE_DECLINE_CHALLENGE:
		handleMsgErr = m.HandleDeclineChallengeMessage(msg)
	case models.CONTENT_TYPE_REVOKE_CHALLENGE:
		handleMsgErr = m.HandleRevokeChallengeMessage(msg)
	}
	if handleMsgErr != nil {
		m.LoggerService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
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
	subErr := m.SubscriptionService.SubClient(msg.SenderKey, msgContent.Topic)
	return subErr
}

func (m *MessageService) HandleEchoMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.EchoMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to EchoMessageContent")
	}
	m.Dispatch(NewEchoEvent(msgContent.Message))
	return nil
}

func (m *MessageService) HandleRequestUpgradeAuthMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.UpgradeAuthRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
	}
	return m.AuthenticationService.UpgradeAuth(msg.SenderKey, msgContent.Role, msgContent.Secret)
}

func (m *MessageService) HandleMoveMessage(moveMsg *models.Message) error {
	moveMsgContent, ok := moveMsg.Content.(*models.MoveMessageContent)
	if !ok {
		return fmt.Errorf("invalid move message content")
	}
	moveErr := m.MatchService.ExecuteMove(moveMsgContent.MatchId, moveMsgContent.Move)
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
	challengeErr := m.MatchService.ChallengePlayer(challengeMsgContent.Challenge)
	if challengeErr != nil {
		return nil
	}
	return nil
}

func (m *MessageService) HandleAcceptChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.AcceptChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid accept challenge message content")
	}
	acceptChallengeErr := m.MatchService.AcceptChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if acceptChallengeErr != nil {
		return nil
	}
	return nil
}

func (m *MessageService) HandleDeclineChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.DeclineChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid decline challenge message content")
	}
	declineChallengeErr := m.MatchService.DeclineChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if declineChallengeErr != nil {
		m.LoggerService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not decline challenge: %s", declineChallengeErr))
	}
	return nil
}

func (m *MessageService) HandleRevokeChallengeMessage(msg *models.Message) error {
	msgContent, ok := msg.Content.(*models.RevokeChallengeMessageContent)
	if !ok {
		return fmt.Errorf("invalid revoke challenge message content")
	}
	revokeChallengeErr := m.MatchService.RevokeChallenge(msg.SenderKey, msgContent.ChallengerClientKey)
	if revokeChallengeErr != nil {
		m.LoggerService.LogRed(models.ENV_SERVER, fmt.Sprintf("could not revoke challenge: %s", revokeChallengeErr))
	}
	return nil
}
