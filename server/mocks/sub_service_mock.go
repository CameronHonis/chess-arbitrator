package mocks

import (
	. "github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/set"
	. "github.com/CameronHonis/stub"
)

type SubServiceMock struct {
	Stubbed[SubscriptionService]
	ServiceMock
}

func NewSubServiceMock(subService *SubscriptionService) *SubServiceMock {
	s := &SubServiceMock{}
	s.Stubbed = *NewStubbed(s, subService)
	s.ServiceMock = *NewServiceMock(&subService.Service)
	return s
}

func (s *SubServiceMock) SubClient(clientKey Key, topic MessageTopic) error {
	out := s.Call("SubClient", clientKey, topic)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (s *SubServiceMock) UnsubClient(clientKey Key, topic MessageTopic) error {
	out := s.Call("UnsubClient", clientKey, topic)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (s *SubServiceMock) UnsubClientFromAll(clientKey Key) {
	_ = s.Call("UnsubClientFromAll", clientKey)
}

func (s *SubServiceMock) SubbedTopics(clientKey Key) *Set[MessageTopic] {
	out := s.Call("SubbedTopics", clientKey)
	return out[0].(*Set[MessageTopic])
}

func (s *SubServiceMock) ClientKeysSubbedToTopic(topic MessageTopic) *Set[Key] {
	out := s.Call("ClientKeysSubbedToTopic", topic)
	return out[0].(*Set[Key])
}
