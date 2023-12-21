package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
)

type MessageConfig struct {
}

func NewMessageConfig() *MessageConfig {
	return &MessageConfig{}
}

func (mhc *MessageConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*MessageConfig))
	return &newConfig
}

type MessageServiceI interface {
	HandleMessage(msg *Message)
}

type MessageService struct {
	Service[*MessageConfig]

	__dependencies__   Marker
	LoggerService      LoggerServiceI
	AuthService        AuthenticationServiceI
	SubService         SubscriptionServiceI
	MatchService       MatchServiceI
	MatchmakingService MatchmakingServiceI

	__state__ Marker
}

func NewMessageHandlerService() *MessageService {
	messageHandler := &MessageService{}
	messageHandler.Service = *NewService(messageHandler, NewMessageConfig())

	return messageHandler
}

func (mh *MessageService) HandleMessage(msg *Message) {
	mh.LoggerService.Log(ENV_CLIENT, fmt.Sprintf("handling message %s", msg))
	var handleMsgErr error
	switch msg.ContentType {
	case CONTENT_TYPE_FIND_MATCH:
		handleMsgErr = mh.HandleFindMatchMessage(msg)
	case CONTENT_TYPE_FIND_BOT_MATCH:
		handleMsgErr = mh.HandleFindBotMatchMessage(msg)
	case CONTENT_TYPE_ECHO:
		handleMsgErr = mh.HandleEchoMessage(msg)
	case CONTENT_TYPE_SUBSCRIBE_REQUEST:
		handleMsgErr = mh.HandleSubscribeRequestMessage(msg)
	case CONTENT_TYPE_UPGRADE_AUTH_REQUEST:
		handleMsgErr = mh.HandleRequestUpgradeAuthMessage(msg)
	case CONTENT_TYPE_INIT_BOT_MATCH_SUCCESS:
		handleMsgErr = mh.HandleInitBotMatchSuccessMessage(msg)
	case CONTENT_TYPE_INIT_BOT_MATCH_FAILURE:
		handleMsgErr = mh.HandleInitBotMatchFailureMessage(msg)
	case CONTENT_TYPE_MOVE:
		handleMsgErr = mh.HandleMoveMessage(msg)
	case CONTENT_TYPE_CHALLENGE_PLAYER:
		handleMsgErr = mh.HandleChallengePlayerMessage(msg)
	case CONTENT_TYPE_CHALLENGE_TERMINATED:
		handleMsgErr = mh.HandleChallengeTerminatedMessage(msg)
	}
	if handleMsgErr != nil {
		mh.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
	}
	GetUserClientsManager().BroadcastMessage(msg)
}

func (mh *MessageService) HandleFindMatchMessage(msg *Message) error {
	// TODO: query for elo, winStreak, lossStreak
	return mh.MatchmakingService.AddClient(&ClientProfile{
		ClientKey:  msg.SenderKey,
		Elo:        1000,
		WinStreak:  0,
		LossStreak: 0,
	})
}

func (mh *MessageService) HandleFindBotMatchMessage(msg *Message) error {
	//msgContent, ok := msg.Content.(*FindBotMatchMessageContent)
	//if !ok {
	//	return fmt.Errorf("could not cast message to FindBotMatchMessageContent")
	//}
	//botClientKey, botKeyErr := mh.authManager.GetBotKey()
	//if botKeyErr != nil {
	//	return fmt.Errorf("could not get bot key: %s", botKeyErr)
	//}
	//if botClientKey == "" {
	//	msg := &Message{
	//		Topic:       "directMessage",
	//		ContentType: CONTENT_TYPE_FIND_BOT_MATCH_NO_BOTS,
	//		Content:     &FindBotMatchNoBotsMessageContent{},
	//	}
	//	return mh.userClientsManager.DirectMessage(msg, msg.SenderKey)
	//}
	//match := NewMatch(msg.SenderKey, botClientKey, NewBulletTimeControl())
	//GetMatchManager().StageMatchFromChallenge(match)
	//outMsg := &Message{
	//	Topic:       "directMessage",
	//	ContentType: CONTENT_TYPE_INIT_BOT_MATCH,
	//	Content: &InitBotMatchMessageContent{
	//		BotName: msgContent.BotName,
	//		MatchId: match.Uuid,
	//	},
	//}
	//return .DirectMessage(outMsg, botClientKey)
	return nil
}

func (mh *MessageService) HandleSubscribeRequestMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*SubscribeRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestMessageContent")
	}
	subErr := mh.SubService.SubClientTo(msg.SenderKey, msgContent.Topic)
	if subErr != nil {
		outMsg := &Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED,
			Content: &SubscribeRequestDeniedMessageContent{
				Reason: subErr.Error(),
				Topic:  msgContent.Topic,
			},
		}
		return mh.userClientsManager.DirectMessage(outMsg, msg.SenderKey)
	}
	subGrantedMsg := &Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED,
		Content: &SubscribeRequestGrantedMessageContent{
			Topic: msg.Topic,
		},
	}
	return mh.userClientsManager.DirectMessage(subGrantedMsg, msg.SenderKey)
}

func (mh *MessageService) HandleEchoMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*EchoMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to EchoMessageContent")
	}
	echoMsg := Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_ECHO,
		Content: &EchoMessageContent{
			Message: msgContent.Message,
		},
	}
	return mh.userClientsManager.DirectMessage(&echoMsg, msg.SenderKey)
}

func (mh *MessageService) HandleRequestUpgradeAuthMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*UpgradeAuthRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
	}
	role, upgradeAuthErr := mh.AuthService.UpgradeAuth(msg.SenderKey, msgContent.Secret)
	if upgradeAuthErr != nil {
		outMsg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_UPGRADE_AUTH_DENIED,
			Content: &UpgradeAuthDeniedMessageContent{
				Reason: "unrecognized secret",
			},
		}
		return mh.userClientsManager.DirectMessage(&outMsg, msg.SenderKey)
	}
	outMsg := Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
		Content: &UpgradeAuthGrantedMessageContent{
			UpgradedToRole: role,
		},
	}
	return mh.userClientsManager.DirectMessage(&outMsg, msg.SenderKey)
}

func (mh *MessageService) HandleInitBotMatchSuccessMessage(msg *Message) error {
	_, ok := msg.Content.(*InitBotMatchSuccessMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to InitBotMatchSuccessMessageContent")
	}
	return nil
}

func (mh *MessageService) HandleInitBotMatchFailureMessage(msg *Message) error {
	_, ok := msg.Content.(*InitBotMatchFailureMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to InitBotMatchFailureMessageContent")
	}
	return nil
}
func (mh *MessageService) HandleMoveMessage(moveMsg *Message) error {
	moveMsgContent, ok := moveMsg.Content.(*MoveMessageContent)
	if !ok {
		return fmt.Errorf("invalid move message content")
	}
	return mh.MatchService.ExecuteMove(moveMsgContent.MatchId, moveMsgContent.Move)
}

func (mh *MessageService) HandleChallengePlayerMessage(challengeMsg *Message) error {
	challengeMsgContent, ok := challengeMsg.Content.(*ChallengePlayerMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge message content")
	}
	_, stageMatchErr := mh.MatchService.StageMatchFromChallenge(challengeMsgContent.Challenge)
	if stageMatchErr != nil {
		return fmt.Errorf("could not stage match for challenge with challenger key %s: %s",
			challengeMsgContent.Challenge.ChallengerKey, stageMatchErr)
	}
	//outMsg := &Message{
	//	Topic:       "directMessage",
	//	ContentType: CONTENT_TYPE_CHALLENGE_PLAYER,
	//	Content: &ChallengePlayerMessageContent{
	//		Challenge: challengeMsgContent.Challenge,
	//	},
	//}
	return nil
}

func (mh *MessageService) HandleChallengeTerminatedMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*ChallengeTerminatedMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge terminated message content")
	}

	challenge := msgContent.Challenge
	terminateChallengeErr := mh.MatchService.TerminateChallenge(challenge)
	if terminateChallengeErr != nil {
		return fmt.Errorf("could not terminate challenge: %s", terminateChallengeErr)
	}
	return mh.userClientsManager.DirectMessage(&Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_CHALLENGE_TERMINATED,
		Content:     &ChallengeTerminatedMessageContent{},
	}, challenge.ChallengerKey)
}
