package server

import (
	"fmt"
	. "github.com/CameronHonis/set"
	"sync"
)

var subscriptionManager *SubscriptionManager

type MessageTopic string
type SubscriptionManagerI interface {
	SubClientTo(clientKey string, topic MessageTopic) error
	UnsubClientFrom(clientKey string, topic MessageTopic) error
	UnsubClientFromAll(clientKey string)
	GetSubbedTopics(clientKey string) *Set[MessageTopic]
	GetClientKeysSubbedToTopic(topic MessageTopic) *Set[string]
}

type SubscriptionManager struct {
	userClientsManager UserClientsManagerI
	authManager        AuthManagerI

	subbedClientKeysByTopic map[MessageTopic]*Set[string]
	subbedTopicsByClientKey map[string]*Set[MessageTopic]
	mu                      sync.Mutex
}

func GetSubscriptionManager() *SubscriptionManager {
	if subscriptionManager != nil {
		return subscriptionManager
	}
	subscriptionManager = &SubscriptionManager{} // null service to prevent infinite recursion
	subscriptionManager = &SubscriptionManager{
		userClientsManager:      GetUserClientsManager(),
		authManager:             GetAuthManager(),
		subbedClientKeysByTopic: make(map[MessageTopic]*Set[string]),
		subbedTopicsByClientKey: make(map[string]*Set[MessageTopic]),
	}
	return subscriptionManager
}

func (sm *SubscriptionManager) SubClientTo(clientKey string, topic MessageTopic) error {
	authErr := sm.authManager.ValidateClientForTopic(clientKey, topic)
	if authErr != nil {
		return authErr
	}
	subbedTopics := sm.GetSubbedTopics(clientKey)
	sm.mu.Lock()
	if subbedTopics.Has(topic) {
		return fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Add(topic)
	sm.mu.Unlock()

	subbedClientKeys := sm.GetClientKeysSubbedToTopic(topic)

	sm.mu.Lock()
	subbedClientKeys.Add(clientKey)
	sm.mu.Unlock()
	return nil
}

func (sm *SubscriptionManager) UnsubClientFrom(clientKey string, topic MessageTopic) error {
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

func (sm *SubscriptionManager) UnsubClientFromAll(clientKey string) {
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

func (sm *SubscriptionManager) GetSubbedTopics(clientKey string) *Set[MessageTopic] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_, ok := sm.subbedTopicsByClientKey[clientKey]
	if !ok {
		sm.subbedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
	}
	return sm.subbedTopicsByClientKey[clientKey]
}

func (sm *SubscriptionManager) GetClientKeysSubbedToTopic(topic MessageTopic) *Set[string] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_, ok := sm.subbedClientKeysByTopic[topic]
	if !ok {
		sm.subbedClientKeysByTopic[topic] = EmptySet[string]()
	}
	return sm.subbedClientKeysByTopic[topic]
}
