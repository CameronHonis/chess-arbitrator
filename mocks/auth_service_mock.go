package mocks

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/stub"
)

type AuthServiceMock struct {
	Stubbed[auth.AuthenticationService]
	ServiceMock
}

func NewAuthServiceMock(authService *auth.AuthenticationService) *AuthServiceMock {
	as := &AuthServiceMock{}
	as.Stubbed = *NewStubbed(as, authService)
	as.ServiceMock = *NewServiceMock(&authService.Service)
	return as
}

func (as *AuthServiceMock) GetRole(clientKey Key) (RoleName, error) {
	out := as.Call("GetRole", clientKey)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(RoleName), err
}

func (as *AuthServiceMock) AddClient(clientKey Key) {
	_ = as.Call("AddClient", clientKey)
}

func (as *AuthServiceMock) UpgradeAuth(clientKey Key, roleName RoleName, secret string) error {
	out := as.Call("UpgradeAuth", clientKey, roleName, secret)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) RemoveClient(clientKey Key) error {
	out := as.Call("RemoveClient", clientKey)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) ValidateSecret(roleName RoleName, secret string) error {
	out := as.Call("ValidateSecret", roleName, secret)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) ValidateAuthInMessage(msg *Message) error {
	out := as.Call("ValidateAuthInMessage", msg)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) StripAuthFromMessage(msg *Message) {
	_ = as.Call("StripAuthFromMessage", msg)
}

func (as *AuthServiceMock) ValidateClientForTopic(clientKey Key, topic MessageTopic) error {
	out := as.Call("ValidateClientForTopic", clientKey, topic)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) SetRole(clientKey Key, role RoleName) error {
	out := as.Call("SetRole", clientKey, role)
	var err error
	if out[0] != nil {
		err = out[0].(error)
	}
	return err
}

func (as *AuthServiceMock) GetSecret(role RoleName) (string, error) {
	out := as.Call("GetSecret", role)
	var err error
	if out[1] != nil {
		err = out[1].(error)
	}
	return out[0].(string), err
}
