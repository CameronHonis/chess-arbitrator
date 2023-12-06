package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/set"
)

type MessageHandlerI interface {
	HandleMessage(msg *Message)
}

var messageHandler *MessageHandler

type MessageHandler struct {
	logManager          LogManagerI
	authManager         AuthManagerI
	matchManager        MatchManagerI
	matchmakingManager  MatchmakingManagerI
	SubscriptionManager SubscriptionManagerI
}

func GetMessageHandler() *MessageHandler {
	if messageHandler == nil {
		messageHandler = &MessageHandler{
			logManager:         GetLogManager(),
			authManager:        GetAuthManager(),
			matchManager:       GetMatchManager(),
			matchmakingManager: GetMatchmakingManager(),
		}
	}
	return messageHandler
}

func (mh *MessageHandler) HandleMessage(msg *Message) {
	mh.logManager.Log(ENV_CLIENT, fmt.Sprintf("handling message %s", msg))
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
		GetLogManager().LogRed(ENV_SERVER, fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
	}
	GetUserClientsManager().BroadcastMessage(msg)
}

func (mh *MessageHandler) HandleFindMatchMessage(msg *Message) error {
	// TODO: query for elo, winStreak, lossStreak
	addClientErr := GetMatchmakingManager().AddClient(&ClientProfile{
		ClientKey:  clientKey,
		Elo:        1000,
		WinStreak:  0,
		LossStreak: 0,
	})
	if addClientErr != nil {
		return fmt.Errorf("could not add client %s to matchmaking pool: %s", clientKey, addClientErr)
	}
	return nil
}

func (mh *MessageHandler) HandleFindBotMatchMessage(msg *Message) error {
	botClientKey := GetAuthManager().chessBotKey
	if botClientKey == "" {
		msg := &Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_FIND_BOT_MATCH_NO_BOTS,
			Content:     &FindBotMatchNoBotsMessageContent{},
		}
		return GetUserClientsManager().DirectMessage(msg, clientKey)
	}
	match := NewMatch(clientKey, botClientKey, NewBulletTimeControl())
	GetMatchManager().StageMatch(match)
	msg := &Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_INIT_BOT_MATCH,
		Content: &InitBotMatchMessageContent{
			BotName: botName,
			MatchId: match.Uuid,
		},
	}
	return GetUserClientsManager().DirectMessage(msg, botClientKey)
}

func (mh *MessageHandler) HandleSubscribeRequestMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*SubscribeRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestMessageContent")
	}
	// TODO: add auth groups - including one for bots client
	topicWhitelist := EmptySet[string]()
	topicWhitelist.Add("findBotMatch")

	if !topicWhitelist.Has(string(topic)) {
		msg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED,
			Content: &SubscribeRequestDeniedMessageContent{
				Topic:  topic,
				Reason: "topic not whitelisted to public",
			},
		}
		return GetUserClientsManager().DirectMessage(&msg, clientKey)
	}
	subErr := GetUserClientsManager().SubscribeClientTo(clientKey, topic)
	if subErr != nil {
		return fmt.Errorf("could not subscribe client %s to topic %s: %s", clientKey, topic, subErr)
	}
	subGrantedMsg := &Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED,
		Content: &SubscribeRequestGrantedMessageContent{
			Topic: topic,
		},
	}
	return GetUserClientsManager().DirectMessage(subGrantedMsg, clientKey)
}

func (mh *MessageHandler) HandleEchoMessage(msg *Message) error {
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
	return GetUserClientsManager().DirectMessage(&echoMsg, msg.SenderKey)
}

func (mh *MessageHandler) HandleRequestUpgradeAuthMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*UpgradeAuthRequestMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
	}
	return mh.authManager.UpgradeAuth(msg.SenderKey, msgContent.Secret)
}

func (mh *MessageHandler) HandleInitBotMatchSuccessMessage(msg *Message) error {
	_, ok := msg.Content.(*InitBotMatchSuccessMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to InitBotMatchSuccessMessageContent")
	}
	return nil
}

func (mh *MessageHandler) HandleInitBotMatchFailureMessage(msg *Message) error {
	_, ok := msg.Content.(*InitBotMatchFailureMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to InitBotMatchFailureMessageContent")
	}
	return nil
}
func (mh *MessageHandler) HandleMoveMessage(moveMsg *Message) error {
	moveMsgContent, ok := moveMsg.Content.(*MoveMessageContent)
	if !ok {
		return fmt.Errorf("invalid move message content")
	}
	return mh.matchManager.ExecuteMove(moveMsgContent.MatchId, moveMsgContent.Move)
}

func (mh *MessageHandler) HandleChallengePlayerMessage(challengeMsg *Message) error {
	challengeMsgContent, ok := challengeMsg.Content.(*ChallengePlayerMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge message content")
	}
	return mh.matchManager.ChallengeClient(challengeMsgContent.Challenge)
}

func (mh *MessageHandler) HandleChallengeTerminatedMessage(msg *Message) error {
	msgContent, ok := msg.Content.(*ChallengeTerminatedMessageContent)
	if !ok {
		return fmt.Errorf("invalid challenge terminated message content")
	}
	return mh.matchManager.TerminateChallenge(msgContent.Challenge, msgContent.Reason)
}
