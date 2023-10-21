package server

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Topic   MessageTopic `json:"topic"`
	Content interface{}  `json:"content"`
}

func (m *Message) IsPrivate() bool {
	return m.Topic == MESSAGE_TOPIC_AUTH
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
		return nil, fmt.Errorf("could not unmarshal json \n%s\n while constructing Message", string(msgJson))
	}
	if msg.Topic == MESSAGE_TOPIC_NONE {
		return nil, fmt.Errorf("could not determine required field 'Topic' from %s while constructing Message", string(msgJson))
	}
	contentMap, isMap := msg.Content.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("could not extract content map from %s while constructing Message content", string(msgJson))
	}
	contentJson, _ := json.Marshal(contentMap)

	var msgContent interface{}
	var contentJsonParseErr error
	switch msg.Topic {
	case MESSAGE_TOPIC_AUTH:
		msgContent = &AuthMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	case MESSAGE_TOPIC_INIT_BOT_MATCH:
		msgContent = &InitBotMatchMessageContent{}
		contentJsonParseErr = json.Unmarshal(contentJson, msgContent)
	default:
		contentJsonParseErr = fmt.Errorf("unhandled content constructor for topic %d while constructing Message.Content", msg.Topic)
	}
	if contentJsonParseErr != nil {
		return nil, contentJsonParseErr
	}

	msg.Content = msgContent
	return &msg, nil
}

type AuthMessageContent struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type InitBotMatchMessageContent struct {
	MatchKey string  `json:"matchKey"`
	BotType  BotType `json:"botType"`
}

type InitMatchMessageContent struct {
	RequesterElo int `json:"requesterElo"`
}
