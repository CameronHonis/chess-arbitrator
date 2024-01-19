package models

import (
	"encoding/json"
	"fmt"
	"github.com/CameronHonis/chess"
)

type MessageTopic string

type Message struct {
	SenderKey   Key          `json:"senderKey"`
	PrivateKey  Key          `json:"privateKey"`
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
		CONTENT_TYPE_FIND_MATCH:                &FindMatchMessageContent{},
		CONTENT_TYPE_MATCH_UPDATE:              &MatchUpdateMessageContent{},
		CONTENT_TYPE_MOVE:                      &MoveMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST:         &SubscribeRequestMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED: &SubscribeRequestGrantedMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:  &SubscribeRequestDeniedMessageContent{},
		CONTENT_TYPE_ECHO:                      &EchoMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_REQUEST:      &UpgradeAuthRequestMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_GRANTED:      &UpgradeAuthGrantedMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_DENIED:       &UpgradeAuthDeniedMessageContent{},
		CONTENT_TYPE_CHALLENGE_REQUEST:         &ChallengePlayerMessageContent{},
		CONTENT_TYPE_CHALLENGE_REQUEST_FAILED:  &ChallengeRequestFailedMessageContent{},
		CONTENT_TYPE_CHALLENGE_ACCEPTED:        &ChallengeAcceptedMessageContent{},
		CONTENT_TYPE_CHALLENGE_DECLINED:        &ChallengeDeclinedMessageContent{},
		CONTENT_TYPE_CHALLENGE_REVOKED:         &ChallengeRevokedMessageContent{},
		CONTENT_TYPE_MATCH_CREATION_FAILED:     &MatchCreationFailedMessageContent{},
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
	CONTENT_TYPE_ECHO                      = "ECHO"
	CONTENT_TYPE_AUTH                      = "AUTH"
	CONTENT_TYPE_FIND_MATCH                = "FIND_MATCH"
	CONTENT_TYPE_MATCH_UPDATE              = "MATCH_UPDATE"
	CONTENT_TYPE_MOVE                      = "MOVE"
	CONTENT_TYPE_MOVE_FAILED               = "MOVE_FAILED"
	CONTENT_TYPE_SUBSCRIBE_REQUEST         = "SUBSCRIBE_REQUEST"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED = "SUBSCRIBE_REQUEST_GRANTED"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED  = "SUBSCRIBE_REQUEST_DENIED"
	CONTENT_TYPE_UPGRADE_AUTH_REQUEST      = "UPGRADE_AUTH_REQUEST"
	CONTENT_TYPE_UPGRADE_AUTH_GRANTED      = "UPGRADE_AUTH_GRANTED"
	CONTENT_TYPE_UPGRADE_AUTH_DENIED       = "UPGRADE_AUTH_DENIED"
	CONTENT_TYPE_CHALLENGE_REQUEST         = "CHALLENGE_REQUEST"
	CONTENT_TYPE_CHALLENGE_REQUEST_FAILED  = "CHALLENGE_REQUEST_FAILED"
	CONTENT_TYPE_CHALLENGE_ACCEPTED        = "ACCEPT_CHALLENGE"
	CONTENT_TYPE_CHALLENGE_DECLINED        = "DECLINE_CHALLENGE"
	CONTENT_TYPE_CHALLENGE_REVOKED         = "REVOKE_CHALLENGE"
	CONTENT_TYPE_MATCH_CREATION_FAILED     = "MATCH_CREATION_FAILED"
)

type AuthMessageContent struct {
	PublicKey  Key `json:"publicKey"`
	PrivateKey Key `json:"privateKey"`
}

type FindMatchMessageContent struct {
}

type MatchUpdateMessageContent struct {
	Match *Match `json:"matcher"`
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
	Role   RoleName `json:"role"`
	Secret string   `json:"secret"`
}

type UpgradeAuthGrantedMessageContent struct {
	UpgradedToRole RoleName `json:"upgradedToRole"`
}

type UpgradeAuthDeniedMessageContent struct {
	Reason string `json:"reason"`
}

type ChallengePlayerMessageContent struct {
	Challenge *Challenge `json:"challenge"`
}

type ChallengeRequestFailedMessageContent struct {
	Challenge *Challenge `json:"challenge"`
	Reason    string     `json:"reason"`
}

type ChallengeAcceptedMessageContent struct {
	ChallengerClientKey Key `json:"challengerClientKey"`
}

type ChallengeDeclinedMessageContent struct {
	ChallengerClientKey Key `json:"challengerClientKey"`
}

type ChallengeRevokedMessageContent struct {
	ChallengerClientKey Key `json:"challengerClientKey"`
}

type MatchCreationFailedMessageContent struct {
	WhiteClientKey Key    `json:"whiteClientKey"`
	BlackClientKey Key    `json:"blackClientKey"`
	Reason         string `json:"reason"`
}
