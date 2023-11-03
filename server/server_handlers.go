package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/set"
)

func HandleMessage(msg *Message, clientKey string) {
	GetLogManager().Log("server", fmt.Sprintf("handling message %s", msg))
	var handleMsgErr error
	switch msg.ContentType {
	case CONTENT_TYPE_FIND_MATCH:
		handleMsgErr = HandleFindMatchMessage(clientKey)
	case CONTENT_TYPE_FIND_BOT_MATCH:
		handleMsgErr = HandleFindBotMatchMessage(clientKey)
	case CONTENT_TYPE_ECHO:
		msgContent, ok := msg.Content.(*EchoMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to EchoMessageContent")
			break
		}
		handleMsgErr = HandleEchoMessage(clientKey, msgContent.Message)
	case CONTENT_TYPE_SUBSCRIBE_REQUEST:
		msgContent, ok := msg.Content.(*SubscribeRequestMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to SubscribeRequestMessageContent")
			break
		}
		handleMsgErr = HandleSubscribeRequestMessage(clientKey, msgContent.Topic)
	case CONTENT_TYPE_UPGRADE_AUTH_REQUEST:
		msgContent, ok := msg.Content.(*UpgradeAuthRequestMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to UpgradeAuthRequestMessageContent")
			break
		}
		handleMsgErr = HandleRequestUpgradeAuthMessage(clientKey, msgContent.Secret)
	case CONTENT_TYPE_INIT_BOT_SUCCESS:
		msgContent, ok := msg.Content.(*InitBotSuccessMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to InitBotSuccessMessageContent")
			break
		}
		handleMsgErr = HandleInitBotSuccessMessage(msgContent)
	case CONTENT_TYPE_INIT_BOT_FAILURE:
		msgContent, ok := msg.Content.(*InitBotFailureMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to InitBotFailureMessageContent")
			break
		}
		handleMsgErr = HandleInitBotMatchFailureMessage(msgContent)
	case CONTENT_TYPE_MOVE:
		msgContent, ok := msg.Content.(*MoveMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to MoveMessageContent")
			break
		}
		handleMsgErr = HandleMoveMessage(clientKey, msgContent.Move)
	}
	if handleMsgErr != nil {
		GetLogManager().LogRed("server", fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
	}
	GetUserClientsManager().BroadcastMessage(msg)
}

func HandleFindMatchMessage(clientKey string) error {
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

func HandleFindBotMatchMessage(clientKey string) error {
	subbedKeys := GetUserClientsManager().GetClientKeysSubscribedToTopic("findBotMatch")
	if len(subbedKeys.Flatten()) == 0 {
		msg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_FIND_BOT_MATCH_NO_BOTS,
			Content:     &FindBotMatchNoBotsMessageContent{},
		}
		return GetUserClientsManager().DirectMessage(&msg, clientKey)
	}
	return nil
}

func HandleSubscribeRequestMessage(clientKey string, topic MessageTopic) error {
	// TODO: add auth groups - including one for bots client
	topicWhitelist := set.EmptySet[string]()
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

func HandleEchoMessage(clientKey string, msg string) error {
	echoMsg := Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_ECHO,
		Content: &EchoMessageContent{
			Message: msg,
		},
	}
	return GetUserClientsManager().DirectMessage(&echoMsg, clientKey)
}

func HandleRequestUpgradeAuthMessage(clientKey string, secret string) error {
	upgradedToRole, upgradeErr := GetAuthManager().UpgradeAuth(clientKey, secret)
	if upgradeErr != nil {
		msg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_UPGRADE_AUTH_DENIED,
			Content: &UpgradeAuthDeniedMessageContent{
				Reason: upgradeErr.Error(),
			},
		}
		return GetUserClientsManager().DirectMessage(&msg, clientKey)
	}
	msg := Message{
		Topic:       "directMessage",
		ContentType: CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
		Content: &UpgradeAuthGrantedMessageContent{
			UpgradedToRole: upgradedToRole,
		},
	}
	return GetUserClientsManager().DirectMessage(&msg, clientKey)
}

func HandleInitBotSuccessMessage(msgContent *InitBotSuccessMessageContent) error {
	match := NewMatch(msgContent.RequesterClientKey, GetAuthManager().chessBotKey, &TimeControl{
		InitialTimeSeconds:  300,
		IncrementSeconds:    0,
		TimeAfterMovesCount: 0,
		SecondsAfterMoves:   0,
	})
	addMatchErr := GetMatchManager().AddMatch(match)
	if addMatchErr != nil {
		return fmt.Errorf("could not init bot match requested by %s: %s", msgContent.RequesterClientKey, addMatchErr)
	}
	// NOTE: probably a bad idea to not establish a topic with a single subscriber (the requester) and broadcast this over the topic
	// 		 but this seemed faster to implement
	return GetUserClientsManager().DirectMessage(
		&Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_INIT_BOT_SUCCESS,
			Content:     msgContent,
		},
		msgContent.RequesterClientKey,
	)
}

func HandleInitBotMatchFailureMessage(msgContent *InitBotFailureMessageContent) error {
	// NOTE: probably a bad idea to not establish a topic with a single subscriber (the requester) and broadcast this over the topic
	// 		 but this seemed faster to implement
	return GetUserClientsManager().DirectMessage(
		&Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_INIT_BOT_FAILURE,
			Content:     msgContent,
		},
		msgContent.RequesterClientKey,
	)
}

func HandleMoveMessage(clientKey string, move *chess.Move) error {
	match, getMatchErr := GetMatchManager().GetMatchByClientId(clientKey)
	if getMatchErr != nil {
		return fmt.Errorf("could not get match for client %s: %s", clientKey, getMatchErr)
	}
	return match.ExecuteMove(move)
}
