package subscription_service

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
)

const (
	SUB_SUCCESSFUL service.EventVariant = "SUB_SUCCESSFUL"
	SUB_FAILED                          = "SUB_FAILED"
)

type SubSuccessfulPayload struct {
	ClientKey models.Key
	Topic     models.MessageTopic
}

type SubSuccessfulEvent struct{ service.Event }

func NewSubSuccessEvent(clientKey models.Key, topic models.MessageTopic) *SubSuccessfulEvent {
	return &SubSuccessfulEvent{
		Event: *service.NewEvent(SUB_SUCCESSFUL, &SubSuccessfulPayload{
			ClientKey: clientKey,
			Topic:     topic,
		}),
	}
}

type SubFailedPayload struct {
	ClientKey models.Key
	Topic     models.MessageTopic
	Reason    string
}

type SubFailedEvent struct{ service.Event }

func NewSubFailedEvent(clientKey models.Key, topic models.MessageTopic, reason string) *SubFailedEvent {
	return &SubFailedEvent{
		Event: *service.NewEvent(SUB_FAILED, &SubFailedPayload{
			ClientKey: clientKey,
			Topic:     topic,
			Reason:    reason,
		}),
	}
}
