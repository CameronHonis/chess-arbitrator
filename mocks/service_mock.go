package mocks

import (
	. "github.com/CameronHonis/service"
	. "github.com/CameronHonis/stub"
)

type ServiceMock struct {
	Stubbed[Service]
}

func NewServiceMock(service *Service) *ServiceMock {
	s := &ServiceMock{}
	s.Stubbed = *NewStubbed[Service](s, service)
	return s
}

func (s *ServiceMock) Config() ConfigI {
	out := s.Call("Config")
	return out[0].(ConfigI)
}

func (s *ServiceMock) Dependencies() []ServiceI {
	out := s.Call("Dependencies")
	return out[0].([]ServiceI)
}

func (s *ServiceMock) AddDependency(service ServiceI) {
	_ = s.Call("AddDependency", service)
}

func (s *ServiceMock) Dispatch(event EventI) {
	_ = s.Call("Dispatch", event)
}

func (s *ServiceMock) AddEventListener(eventVariant EventVariant, fn EventHandler) (eventId int) {
	out := s.Call("AddEventListener", eventVariant, fn)
	return out[0].(int)
}

func (s *ServiceMock) RemoveEventListener(eventId int) {
	_ = s.Call("RemoveEventListener", eventId)
}

func (s *ServiceMock) PropagateEvent(event EventI) {
	_ = s.Call("PropagateEvent", event)
}

func (s *ServiceMock) SetParent(parent ServiceI) {
	_ = s.Call("SetParent", parent)
}
