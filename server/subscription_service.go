package server

import (
	"fmt"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/CameronHonis/set"
	"sync"
)

const (
	SUB_SUCCESSFUL EventVariant = "SUB_SUCCESSFUL"
	SUB_FAILED                  = "SUB_FAILED"
)

type SubSuccessfulPayload struct {
	ClientKey Key
	Topic     MessageTopic
}

type SubSuccessfulEvent struct{ Event }

func NewSubSuccessEvent(clientKey Key, topic MessageTopic) *SubSuccessfulEvent {
	return &SubSuccessfulEvent{
		Event: *NewEvent(SUB_SUCCESSFUL, &SubSuccessfulPayload{
			ClientKey: clientKey,
			Topic:     topic,
		}),
	}
}

type SubFailedPayload struct {
	ClientKey Key
	Topic     MessageTopic
	Reason    string
}

type SubFailedEvent struct{ Event }

func NewSubFailedEvent(clientKey Key, topic MessageTopic, reason string) *SubFailedEvent {
	return &SubFailedEvent{
		Event: *NewEvent(SUB_FAILED, &SubFailedPayload{
			ClientKey: clientKey,
			Topic:     topic,
			Reason:    reason,
		}),
	}
}

type MessageTopic string

type SubscriptionConfig struct {
	ConfigI
}

func NewSubscriptionConfig() *SubscriptionConfig {
	return &SubscriptionConfig{}
}

type SubscriptionServiceI interface {
	SubClientTo(clientKey Key, topic MessageTopic) error
	UnsubClientFrom(clientKey Key, topic MessageTopic) error
	UnsubClientFromAll(clientKey Key)
	GetSubbedTopics(clientKey Key) *Set[MessageTopic]
	GetClientKeysSubbedToTopic(topic MessageTopic) *Set[Key]
}

type SubscriptionService struct {
	Service
	__dependencies__      Marker
	AuthenticationService AuthenticationServiceI

	__state__               Marker
	subbedClientKeysByTopic map[MessageTopic]*Set[Key]
	subbedTopicsByClientKey map[Key]*Set[MessageTopic]
	mu                      sync.Mutex
}

func NewSubscriptionService(config *SubscriptionConfig) *SubscriptionService {
	subService := &SubscriptionService{
		subbedClientKeysByTopic: make(map[MessageTopic]*Set[Key]),
		subbedTopicsByClientKey: make(map[Key]*Set[MessageTopic]),
	}
	subService.Service = *NewService(subService, config)
	return subService
}
func (sm *SubscriptionService) SubClientTo(clientKey Key, topic MessageTopic) error {
	authErr := sm.AuthenticationService.ValidateClientForTopic(clientKey, topic)
	if authErr != nil {
		go sm.Dispatch(NewSubFailedEvent(clientKey, topic, authErr.Error()))
		return authErr
	}
	subbedTopics := sm.GetSubbedTopics(clientKey)
	sm.mu.Lock()
	if subbedTopics.Has(topic) {
		go sm.Dispatch(NewSubFailedEvent(clientKey, topic, "already subscribed"))
		return fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Add(topic)
	sm.mu.Unlock()

	subbedClientKeys := sm.GetClientKeysSubbedToTopic(topic)

	sm.mu.Lock()
	subbedClientKeys.Add(clientKey)
	sm.mu.Unlock()

	sm.Dispatch(NewSubSuccessEvent(clientKey, topic))
	return nil
}

func (sm *SubscriptionService) UnsubClientFrom(clientKey Key, topic MessageTopic) error {
	subbedTopics := sm.GetSubbedTopics(clientKey)
	sm.mu.Lock()
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	sm.mu.Unlock()

	subbedClientKeys := sm.GetClientKeysSubbedToTopic(topic)

	sm.mu.Lock()
	subbedClientKeys.Remove(clientKey)
	sm.mu.Unlock()
	return nil
}

func (sm *SubscriptionService) UnsubClientFromAll(clientKey Key) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	subbedTopics, ok := sm.subbedTopicsByClientKey[clientKey]
	if !ok {
		return
	}
	for _, topic := range subbedTopics.Flatten() {
		sm.subbedClientKeysByTopic[topic].Remove(clientKey)
	}
	delete(sm.subbedTopicsByClientKey, clientKey)
}

func (sm *SubscriptionService) GetSubbedTopics(clientKey Key) *Set[MessageTopic] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_, ok := sm.subbedTopicsByClientKey[clientKey]
	if !ok {
		sm.subbedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
	}
	return sm.subbedTopicsByClientKey[clientKey]
}

func (sm *SubscriptionService) GetClientKeysSubbedToTopic(topic MessageTopic) *Set[Key] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_, ok := sm.subbedClientKeysByTopic[topic]
	if !ok {
		sm.subbedClientKeysByTopic[topic] = EmptySet[Key]()
	}
	return sm.subbedClientKeysByTopic[topic]
}
