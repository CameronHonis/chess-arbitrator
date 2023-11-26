package server

import (
	"fmt"
	"github.com/CameronHonis/chess"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/set"
)

func HandleMessage(msg *Message, clientKey string) {
	GetLogManager().Log(ENV_CLIENT, fmt.Sprintf("handling message %s", msg))
	var handleMsgErr error
	switch msg.ContentType {
	case CONTENT_TYPE_FIND_MATCH:
		handleMsgErr = HandleFindMatchMessage(clientKey)
	case CONTENT_TYPE_FIND_BOT_MATCH:
		msgContent, ok := msg.Content.(*FindBotMatchMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to FindBotMatchMessageContent")
			break
		}
		handleMsgErr = HandleFindBotMatchMessage(clientKey, msgContent.BotName)
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
	case CONTENT_TYPE_INIT_BOT_MATCH_SUCCESS:
		msgContent, ok := msg.Content.(*InitBotMatchSuccessMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to InitBotMatchSuccessMessageContent")
			break
		}
		handleMsgErr = HandleInitBotMatchSuccessMessage(msgContent)
	case CONTENT_TYPE_INIT_BOT_MATCH_FAILURE:
		msgContent, ok := msg.Content.(*InitBotMatchFailureMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to InitBotMatchFailureMessageContent")
			break
		}
		handleMsgErr = HandleInitBotMatchFailureMessage(msgContent)
	case CONTENT_TYPE_MOVE:
		msgContent, ok := msg.Content.(*MoveMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to MoveMessageContent")
			break
		}
		handleMsgErr = HandleMoveMessage(msgContent.MatchId, msgContent.Move)
	case CONTENT_TYPE_CHALLENGE_PLAYER:
		msgContent, ok := msg.Content.(*ChallengePlayerMessageContent)
		if !ok {
			handleMsgErr = fmt.Errorf("could not cast message to ChallengePlayerMessageContent")
			break
		}
		handleMsgErr = HandleChallengePlayerMessage(msgContent)
	}
	if handleMsgErr != nil {
		GetLogManager().LogRed(ENV_SERVER, fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
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

func HandleFindBotMatchMessage(clientKey string, botName string) error {
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

func HandleSubscribeRequestMessage(clientKey string, topic MessageTopic) error {
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

func HandleInitBotMatchSuccessMessage(msgContent *InitBotMatchSuccessMessageContent) error {
	return GetMatchManager().AddMatchFromStaged(msgContent.MatchId)
}

func HandleInitBotMatchFailureMessage(msgContent *InitBotMatchFailureMessageContent) error {
	stagedMatch, fetchStagedMatchErr := GetMatchManager().GetStagedMatchById(msgContent.MatchId)
	if fetchStagedMatchErr != nil {
		return fmt.Errorf("could not fetch staged match with id %s: %s", msgContent.MatchId, fetchStagedMatchErr)
	}
	var clientKey string
	if stagedMatch.WhiteClientId != GetAuthManager().chessBotKey {
		clientKey = stagedMatch.WhiteClientId
	} else {
		clientKey = stagedMatch.BlackClientId
	}

	GetMatchManager().UnstageMatch(msgContent.MatchId)

	return GetUserClientsManager().DirectMessage(
		&Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_INIT_BOT_MATCH_FAILURE,
			Content:     msgContent,
		},
		clientKey,
	)
}

func HandleMoveMessage(matchId string, move *chess.Move) error {
	return GetMatchManager().ExecuteMove(matchId, move)
}

func HandleChallengePlayerMessage(msgContent *ChallengePlayerMessageContent) error {
	return GetMatchManager().ChallengeClient(msgContent.Challenge)
}
