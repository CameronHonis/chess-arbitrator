package subscription_service

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/CameronHonis/set"
	"sync"
)

type SubscriptionServiceI interface {
	ServiceI
	SubClient(clientKey Key, topic MessageTopic) error
	UnsubClient(clientKey Key, topic MessageTopic) error
	UnsubClientFromAll(clientKey Key)
	SubbedTopics(clientKey Key) *Set[MessageTopic]
	ClientKeysSubbedToTopic(topic MessageTopic) *Set[Key]
}

type SubscriptionService struct {
	Service
	__dependencies__ Marker
	AuthService      auth_service.AuthenticationServiceI

	__state__               Marker
	subbedClientKeysByTopic map[MessageTopic]*Set[Key]
	subbedTopicsByClientKey map[Key]*Set[MessageTopic]
	mu                      sync.Mutex
}

func NewSubscriptionService(config *SubscriptionServiceConfig) *SubscriptionService {
	subService := &SubscriptionService{
		subbedClientKeysByTopic: make(map[MessageTopic]*Set[Key]),
		subbedTopicsByClientKey: make(map[Key]*Set[MessageTopic]),
	}
	subService.Service = *NewService(subService, config)
	return subService
}
func (s *SubscriptionService) SubClient(clientKey Key, topic MessageTopic) error {
	authErr := s.AuthService.ValidateClientForTopic(clientKey, topic)
	if authErr != nil {
		go s.Dispatch(NewSubFailedEvent(clientKey, topic, authErr.Error()))
		return authErr
	}
	subbedTopics := s.SubbedTopics(clientKey)
	s.mu.Lock()
	if subbedTopics.Has(topic) {
		go s.Dispatch(NewSubFailedEvent(clientKey, topic, "already subscribed"))
		return fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Add(topic)
	s.mu.Unlock()

	subbedClientKeys := s.ClientKeysSubbedToTopic(topic)

	s.mu.Lock()
	subbedClientKeys.Add(clientKey)
	s.mu.Unlock()

	go s.Dispatch(NewSubSuccessEvent(clientKey, topic))
	return nil
}

func (s *SubscriptionService) UnsubClient(clientKey Key, topic MessageTopic) error {
	subbedTopics := s.SubbedTopics(clientKey)
	s.mu.Lock()
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	s.mu.Unlock()

	subbedClientKeys := s.ClientKeysSubbedToTopic(topic)

	s.mu.Lock()
	subbedClientKeys.Remove(clientKey)
	s.mu.Unlock()
	return nil
}

func (s *SubscriptionService) UnsubClientFromAll(clientKey Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subbedTopics, ok := s.subbedTopicsByClientKey[clientKey]
	if !ok {
		return
	}
	for _, topic := range subbedTopics.Flatten() {
		s.subbedClientKeysByTopic[topic].Remove(clientKey)
	}
	delete(s.subbedTopicsByClientKey, clientKey)
}

func (s *SubscriptionService) SubbedTopics(clientKey Key) *Set[MessageTopic] {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.subbedTopicsByClientKey[clientKey]
	if !ok {
		s.subbedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
	}
	return s.subbedTopicsByClientKey[clientKey]
}

func (s *SubscriptionService) ClientKeysSubbedToTopic(topic MessageTopic) *Set[Key] {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.subbedClientKeysByTopic[topic]
	if !ok {
		s.subbedClientKeysByTopic[topic] = EmptySet[Key]()
	}
	return s.subbedClientKeysByTopic[topic]
}
