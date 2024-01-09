package message_service

import (
	"github.com/CameronHonis/chess-arbitrator/mocks"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/stub"
)

type MessageServiceMock struct {
	Stubbed[MessageService]
	mocks.ServiceMock
}

func NewMessageServiceMock(msgService *MessageService) *MessageServiceMock {
	s := &MessageServiceMock{}
	s.Stubbed = *NewStubbed(s, msgService)
	s.ServiceMock = *mocks.NewServiceMock(&msgService.Service)
	return s
}

func (s *MessageServiceMock) HandleMessage(msg *Message) {
	_ = s.Call("HandleMessage", msg)
}
