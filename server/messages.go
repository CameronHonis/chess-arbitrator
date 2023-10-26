package server

import (
	"encoding/json"
	"fmt"
	"github.com/CameronHonis/chess"
)

type MessageTopic string

type Message struct {
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
	}
	contentMap, isMap := msg.Content.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("could not extract content map from %s while constructing Message content", string(msgJson))
	}
	contentJson, _ := json.Marshal(contentMap)

	var msgContent interface{}
	var contentJsonParseErr error
	switch msg.ContentType {
	case CONTENT_TYPE_EMPTY:
		if msg.Content != nil {
			return nil, fmt.Errorf("message with content type %s has non-nil content", CONTENT_TYPE_EMPTY)
		}
	}
	switch msg.ContentType {
	case CONTENT_TYPE_AUTH:
		msgContent = &AuthMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	case CONTENT_TYPE_FIND_BOT_MATCH:
		msgContent = &FindBotMatchMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	case CONTENT_TYPE_FIND_MATCH:
		msgContent = &FindMatchMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	case CONTENT_TYPE_MATCH_UPDATE:
		msgContent = &MatchUpdateMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	case CONTENT_TYPE_MOVE:
		msgContent = &MoveMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	default:
		contentJsonParseErr = fmt.Errorf("could not extract content type from %s while constructing Message content", msgJson)
	}
	if contentJsonParseErr != nil {
		return nil, contentJsonParseErr
	}

	msg.Content = msgContent
	return &msg, nil
}

type ContentType string

const (
	CONTENT_TYPE_EMPTY          = "EMPTY"
	CONTENT_TYPE_AUTH           = "AUTH"
	CONTENT_TYPE_FIND_BOT_MATCH = "FIND_BOT_MATCH"
	CONTENT_TYPE_FIND_MATCH     = "FIND_MATCH"
	CONTENT_TYPE_MATCH_UPDATE   = "MATCH_UPDATE"
	CONTENT_TYPE_MOVE           = "MOVE"
)

type AuthMessageContent struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type FindBotMatchMessageContent struct {
	MatchKey string `json:"matchKey"`
	BotName  string `json:"botType"`
}

type FindMatchMessageContent struct {
}

type MatchUpdateMessageContent struct {
	Match *Match `json:"match"`
}

type MoveMessageContent struct {
	Move chess.Move `json:"move"`
}
