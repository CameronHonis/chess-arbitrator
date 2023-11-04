package server

import (
	"encoding/json"
	"fmt"
	"github.com/CameronHonis/chess"
)

type MessageTopic string

type Message struct {
	SenderKey   string       `json:"senderKey"`
	PrivateKey  string       `json:"privateKey"`
	Topic       MessageTopic `json:"topic"`
	ContentType ContentType  `json:"contentType"`
	Content     interface{}  `json:"content"`
}

func (m *Message) Marshal() ([]byte, error) {
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func UnmarshalToMessage(msgJson []byte) (*Message, error) {
	var msg Message
	jsonParseErr := json.Unmarshal(msgJson, &msg)
	if jsonParseErr != nil {
		return nil, fmt.Errorf("could not unmarshal json... \n\n\t     %s \n\n\t ...while constructing Message", string(msgJson))
	}

	if msg.ContentType == CONTENT_TYPE_EMPTY && msg.Content != nil {
		return nil, fmt.Errorf("message with content type %s has non-nil content", CONTENT_TYPE_EMPTY)
	} else if msg.ContentType == "" {
		return nil, fmt.Errorf("could not extract content type from %s while constructing Message content", string(msgJson))
	}
	contentMap, isMap := msg.Content.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("could not extract content map from %s while constructing Message content", string(msgJson))
	}
	contentJson, _ := json.Marshal(contentMap)

	var contentJsonParseErr error
	msg.Content, contentJsonParseErr = UnmarshalMessageContent(msg.ContentType, contentJson)
	if contentJsonParseErr != nil {
		return nil, contentJsonParseErr
	}
	return &msg, nil
}

func UnmarshalMessageContent(contentType ContentType, contentJson []byte) (interface{}, error) {
	contentStructMap := map[ContentType]interface{}{
		CONTENT_TYPE_AUTH:                      &AuthMessageContent{},
		CONTENT_TYPE_FIND_BOT_MATCH:            &FindBotMatchMessageContent{},
		CONTENT_TYPE_FIND_MATCH:                &FindMatchMessageContent{},
		CONTENT_TYPE_MATCH_UPDATE:              &MatchUpdateMessageContent{},
		CONTENT_TYPE_MOVE:                      &MoveMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST:         &SubscribeRequestMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED: &SubscribeRequestGrantedMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:  &SubscribeRequestDeniedMessageContent{},
		CONTENT_TYPE_FIND_BOT_MATCH_NO_BOTS:    &FindBotMatchNoBotsMessageContent{},
		CONTENT_TYPE_ECHO:                      &EchoMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_REQUEST:      &UpgradeAuthRequestMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_GRANTED:      &UpgradeAuthGrantedMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_DENIED:       &UpgradeAuthDeniedMessageContent{},
		CONTENT_TYPE_INIT_BOT_MATCH:            &InitBotMatchMessageContent{},
		CONTENT_TYPE_INIT_BOT_MATCH_SUCCESS:    &InitBotMatchSuccessMessageContent{},
		CONTENT_TYPE_INIT_BOT_MATCH_FAILURE:    &InitBotMatchFailureMessageContent{},
	}
	msgContent, ok := contentStructMap[contentType]
	if !ok {
		return nil, fmt.Errorf("contentStructMap does not specify map between content type %s and existing struct", contentType)
	}
	contentJsonParseErr := json.Unmarshal(contentJson, msgContent)
	if contentJsonParseErr != nil {
		return nil, contentJsonParseErr
	}
	return msgContent, nil
}

type ContentType string

const (
	CONTENT_TYPE_EMPTY                     = "EMPTY"
	CONTENT_TYPE_AUTH                      = "AUTH"
	CONTENT_TYPE_FIND_BOT_MATCH            = "FIND_BOT_MATCH"
	CONTENT_TYPE_FIND_MATCH                = "FIND_MATCH"
	CONTENT_TYPE_MATCH_UPDATE              = "MATCH_UPDATE"
	CONTENT_TYPE_MOVE                      = "MOVE"
	CONTENT_TYPE_SUBSCRIBE_REQUEST         = "SUBSCRIBE_REQUEST"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED = "SUBSCRIBE_REQUEST_GRANTED"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED  = "SUBSCRIBE_REQUEST_DENIED"
	CONTENT_TYPE_FIND_BOT_MATCH_NO_BOTS    = "FIND_BOT_MATCH_NO_BOTS"
	CONTENT_TYPE_ECHO                      = "ECHO"
	CONTENT_TYPE_UPGRADE_AUTH_REQUEST      = "UPGRADE_AUTH_REQUEST"
	CONTENT_TYPE_UPGRADE_AUTH_GRANTED      = "UPGRADE_AUTH_GRANTED"
	CONTENT_TYPE_UPGRADE_AUTH_DENIED       = "UPGRADE_AUTH_DENIED"
	CONTENT_TYPE_INIT_BOT_MATCH            = "INIT_BOT_MATCH"
	CONTENT_TYPE_INIT_BOT_MATCH_SUCCESS    = "INIT_BOT_MATCH_SUCCESS"
	CONTENT_TYPE_INIT_BOT_MATCH_FAILURE    = "INIT_BOT_MATCH_FAILURE"
)

type AuthMessageContent struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type FindBotMatchMessageContent struct {
	BotName string `json:"botName"`
}

type FindBotMatchNoBotsMessageContent struct {
}

type FindMatchMessageContent struct {
}

type MatchUpdateMessageContent struct {
	Match *Match `json:"match"`
}

type MoveMessageContent struct {
	MatchId string      `json:"matchId"`
	Move    *chess.Move `json:"move"`
}

type SubscribeRequestMessageContent struct {
	Topic MessageTopic `json:"topic"`
}

type SubscribeRequestGrantedMessageContent struct {
	Topic MessageTopic `json:"topic"`
}

type SubscribeRequestDeniedMessageContent struct {
	Topic  MessageTopic `json:"topic"`
	Reason string       `json:"reason"`
}

type EchoMessageContent struct {
	Message string `json:"message"`
}

type UpgradeAuthRequestMessageContent struct {
	Secret string `json:"secret"`
}

type UpgradeAuthGrantedMessageContent struct {
	UpgradedToRole string `json:"upgradedToRole"`
}

type UpgradeAuthDeniedMessageContent struct {
	Reason string `json:"reason"`
}

type InitBotMatchMessageContent struct {
	BotName string `json:"botName"`
	MatchId string `json:"matchId"`
}

type InitBotMatchSuccessMessageContent struct {
	BotName string `json:"botType"`
	MatchId string `json:"matchId"`
}

type InitBotMatchFailureMessageContent struct {
	BotName string `json:"botType"`
	MatchId string `json:"matchId"`
	Reason  string `json:"reason"`
}
