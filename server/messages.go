package server

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Topic   MessageTopic `json:"topic"`
	Content interface{}  `json:"content"`
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

	var msgContent interface{}
	var contentCreationErr error
	switch msg.Topic {
	case MESSAGE_TOPIC_AUTH:
		msgContent, contentCreationErr = AuthMessageContentFromMap(contentMap)
	default:
		contentCreationErr = fmt.Errorf("unhandled content constructor for topic %d while constructing Message.Content", msg.Topic)
	}

	if contentCreationErr != nil {
		return nil, contentCreationErr
	}
	msg.Content = msgContent
	return &msg, nil
}

type AuthMessageContent struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func (content *AuthMessageContent) toJsonBytes() ([]byte, error) {
	msgJson, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	return msgJson, nil
}

func AuthMessageContentFromMap(m map[string]interface{}) (*AuthMessageContent, error) {
	pubKey, ok := m["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse publicKey")
	}
	privKey, ok := m["privateKey"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse privateKey")
	}
	return &AuthMessageContent{
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}, nil
}
