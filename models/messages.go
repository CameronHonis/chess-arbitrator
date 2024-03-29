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
		CONTENT_TYPE_REFRESH_AUTH:              &RefreshAuthMessageContent{},
		CONTENT_TYPE_INVALID_AUTH:              &NoMessageContent{},
		CONTENT_TYPE_JOIN_MATCHMAKING:          &FindMatchMessageContent{},
		CONTENT_TYPE_LEAVE_MATCHMAKING:         &NoMessageContent{},
		CONTENT_TYPE_MATCH_UPDATED:             &MatchUpdateMessageContent{},
		CONTENT_TYPE_MOVE:                      &MoveMessageContent{},
		CONTENT_TYPE_RESIGN_MATCH:              &ResignMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST:         &SubscribeRequestMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED: &SubscribeRequestGrantedMessageContent{},
		CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:  &SubscribeRequestDeniedMessageContent{},
		CONTENT_TYPE_ECHO:                      &EchoMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_REQUEST:      &UpgradeAuthRequestMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_GRANTED:      &UpgradeAuthGrantedMessageContent{},
		CONTENT_TYPE_UPGRADE_AUTH_DENIED:       &UpgradeAuthDeniedMessageContent{},
		CONTENT_TYPE_CHALLENGE_REQUEST:         &ChallengeRequestMessageContent{},
		CONTENT_TYPE_CHALLENGE_REQUEST_FAILED:  &ChallengeRequestFailedMessageContent{},
		CONTENT_TYPE_ACCEPT_CHALLENGE:          &AcceptChallengeMessageContent{},
		CONTENT_TYPE_DECLINE_CHALLENGE:         &DeclineChallengeMessageContent{},
		CONTENT_TYPE_REVOKE_CHALLENGE:          &RevokeChallengeMessageContent{},
		CONTENT_TYPE_CHALLENGE_UPDATED:         &ChallengeUpdatedMessageContent{},
		CONTENT_TYPE_MATCH_CREATION_FAILED:     &MatchCreationFailedMessageContent{},
		CONTENT_TYPE_MOVE_FAILED:               &MoveFailedMessageContent{},
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
	// arbitrator responses
	CONTENT_TYPE_AUTH                      ContentType = "AUTH"
	CONTENT_TYPE_INVALID_AUTH              ContentType = "INVALID_AUTH"
	CONTENT_TYPE_MATCH_UPDATED             ContentType = "MATCH_UPDATED"
	CONTENT_TYPE_CHALLENGE_UPDATED         ContentType = "CHALLENGE_UPDATED"
	CONTENT_TYPE_MOVE_FAILED               ContentType = "MOVE_FAILED"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED ContentType = "SUBSCRIBE_REQUEST_GRANTED"
	CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED  ContentType = "SUBSCRIBE_REQUEST_DENIED"
	CONTENT_TYPE_UPGRADE_AUTH_GRANTED      ContentType = "UPGRADE_AUTH_GRANTED"
	CONTENT_TYPE_UPGRADE_AUTH_DENIED       ContentType = "UPGRADE_AUTH_DENIED"
	CONTENT_TYPE_CHALLENGE_REQUEST_FAILED  ContentType = "CHALLENGE_REQUEST_FAILED"
	CONTENT_TYPE_MATCH_CREATION_FAILED     ContentType = "MATCH_CREATION_FAILED"

	// client requests
	CONTENT_TYPE_REFRESH_AUTH         ContentType = "REFRESH_AUTH"
	CONTENT_TYPE_EMPTY                ContentType = "EMPTY"
	CONTENT_TYPE_ECHO                 ContentType = "ECHO"
	CONTENT_TYPE_JOIN_MATCHMAKING     ContentType = "JOIN_MATCHMAKING"
	CONTENT_TYPE_LEAVE_MATCHMAKING    ContentType = "LEAVE_MATCHMAKING"
	CONTENT_TYPE_MOVE                 ContentType = "MOVE"
	CONTENT_TYPE_RESIGN_MATCH         ContentType = "RESIGN_MATCH"
	CONTENT_TYPE_SUBSCRIBE_REQUEST    ContentType = "SUBSCRIBE_REQUEST"
	CONTENT_TYPE_UPGRADE_AUTH_REQUEST ContentType = "UPGRADE_AUTH_REQUEST"
	CONTENT_TYPE_CHALLENGE_REQUEST    ContentType = "CHALLENGE_REQUEST"
	CONTENT_TYPE_ACCEPT_CHALLENGE     ContentType = "ACCEPT_CHALLENGE"
	CONTENT_TYPE_DECLINE_CHALLENGE    ContentType = "DECLINE_CHALLENGE"
	CONTENT_TYPE_REVOKE_CHALLENGE     ContentType = "REVOKE_CHALLENGE"
)

type NoMessageContent struct{}

type AuthMessageContent struct {
	PublicKey  Key `json:"publicKey"`
	PrivateKey Key `json:"privateKey"`
}

type RefreshAuthMessageContent struct {
	ExistingAuth *AuthMessageContent `json:"existingAuth"`
}

type FindMatchMessageContent struct {
	TimeControl *TimeControl `json:"timeControl"`
}

type MatchUpdateMessageContent struct {
	Match *Match `json:"match"`
}

type MoveMessageContent struct {
	MatchId string      `json:"matchId"`
	Move    *chess.Move `json:"move"`
}

type ResignMessageContent struct {
	MatchId string `json:"matchId"`
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

type ChallengeRequestMessageContent struct {
	Challenge *Challenge `json:"challenge"`
}

type ChallengeRequestFailedMessageContent struct {
	Challenge *Challenge `json:"challenge"`
	Reason    string     `json:"reason"`
}

type AcceptChallengeMessageContent struct {
	ChallengerClientKey Key `json:"challengerClientKey"`
}

type DeclineChallengeMessageContent struct {
	ChallengerClientKey Key `json:"challengerClientKey"`
}

type RevokeChallengeMessageContent struct {
	ChallengedClientKey Key `json:"challengedClientKey"`
}

type ChallengeUpdatedMessageContent struct {
	Challenge *Challenge `json:"challenge"`
}

type MatchCreationFailedMessageContent struct {
	WhiteClientKey Key    `json:"whiteClientKey"`
	BlackClientKey Key    `json:"blackClientKey"`
	Reason         string `json:"reason"`
}

type MoveFailedMessageContent struct {
	Move   *chess.Move `json:"move"`
	Reason string      `json:"reason"`
}
