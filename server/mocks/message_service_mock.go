package mocks

import (
	. "github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/stub"
)

type MessageServiceMock struct {
	Stubbed[MessageService]
	ServiceMock
}

func NewMessageServiceMock(msgService *MessageService) *MessageServiceMock {
	s := &MessageServiceMock{}
	s.Stubbed = *NewStubbed(s, msgService)
	s.ServiceMock = *NewServiceMock(&msgService.Service)
	return s
}

func (s *MessageServiceMock) HandleMessage(msg *Message) {
	_ = s.Call("HandleMessage", msg)
}
