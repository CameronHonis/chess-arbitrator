// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CameronHonis/chess-arbitrator/auth (interfaces: AuthenticationServiceI)
//
// Generated by this command:
//
//	mockgen -destination mock/auth_serivce_mock . AuthenticationServiceI
//

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	reflect "reflect"
	models "github.com/CameronHonis/chess-arbitrator/models"
	service "github.com/CameronHonis/service"
	set "github.com/CameronHonis/set"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthenticationServiceI is a mock of AuthenticationServiceI interface.
type MockAuthenticationServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationServiceIMockRecorder
}

// MockAuthenticationServiceIMockRecorder is the mock recorder for MockAuthenticationServiceI.
type MockAuthenticationServiceIMockRecorder struct {
	mock *MockAuthenticationServiceI
}

// NewMockAuthenticationServiceI creates a new mock instance.
func NewMockAuthenticationServiceI(ctrl *gomock.Controller) *MockAuthenticationServiceI {
	mock := &MockAuthenticationServiceI{ctrl: ctrl}
	mock.recorder = &MockAuthenticationServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticationServiceI) EXPECT() *MockAuthenticationServiceIMockRecorder {
	return m.recorder
}

// AddClient mocks base method.
func (m *MockAuthenticationServiceI) AddClient(arg0 models.Key) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClient", arg0)
}

// AddClient indicates an expected call of AddClient.
func (mr *MockAuthenticationServiceIMockRecorder) AddClient(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClient", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddClient), arg0)
}

// AddDependency mocks base method.
func (m *MockAuthenticationServiceI) AddDependency(arg0 service.ServiceI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDependency", arg0)
}

// AddDependency indicates an expected call of AddDependency.
func (mr *MockAuthenticationServiceIMockRecorder) AddDependency(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDependency", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddDependency), arg0)
}

// AddEventListener mocks base method.
func (m *MockAuthenticationServiceI) AddEventListener(arg0 service.EventVariant, arg1 service.EventHandler) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEventListener", arg0, arg1)
	ret0, _ := ret[0].(int)
	return ret0
}

// AddEventListener indicates an expected call of AddEventListener.
func (mr *MockAuthenticationServiceIMockRecorder) AddEventListener(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventListener", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddEventListener), arg0, arg1)
}

// BotClientExists mocks base method.
func (m *MockAuthenticationServiceI) BotClientExists() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BotClientExists")
	ret0, _ := ret[0].(bool)
	return ret0
}

// BotClientExists indicates an expected call of BotClientExists.
func (mr *MockAuthenticationServiceIMockRecorder) BotClientExists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BotClientExists", reflect.TypeOf((*MockAuthenticationServiceI)(nil).BotClientExists))
}

// ClientKeysByRole mocks base method.
func (m *MockAuthenticationServiceI) ClientKeysByRole(arg0 models.RoleName) *set.Set[models.Key] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientKeysByRole", arg0)
	ret0, _ := ret[0].(*set.Set[models.Key])
	return ret0
}

// ClientKeysByRole indicates an expected call of ClientKeysByRole.
func (mr *MockAuthenticationServiceIMockRecorder) ClientKeysByRole(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientKeysByRole", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ClientKeysByRole), arg0)
}

// Config mocks base method.
func (m *MockAuthenticationServiceI) Config() service.ConfigI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(service.ConfigI)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockAuthenticationServiceIMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Config))
}

// Dependencies mocks base method.
func (m *MockAuthenticationServiceI) Dependencies() []service.ServiceI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dependencies")
	ret0, _ := ret[0].([]service.ServiceI)
	return ret0
}

// Dependencies indicates an expected call of Dependencies.
func (mr *MockAuthenticationServiceIMockRecorder) Dependencies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dependencies", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Dependencies))
}

// Dispatch mocks base method.
func (m *MockAuthenticationServiceI) Dispatch(arg0 service.EventI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", arg0)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockAuthenticationServiceIMockRecorder) Dispatch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Dispatch), arg0)
}

// GetRole mocks base method.
func (m *MockAuthenticationServiceI) GetRole(arg0 models.Key) (models.RoleName, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0)
	ret0, _ := ret[0].(models.RoleName)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockAuthenticationServiceIMockRecorder) GetRole(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockAuthenticationServiceI)(nil).GetRole), arg0)
}

// OnStart mocks base method.
func (m *MockAuthenticationServiceI) OnStart() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnStart")
}

// OnStart indicates an expected call of OnStart.
func (mr *MockAuthenticationServiceIMockRecorder) OnStart() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnStart", reflect.TypeOf((*MockAuthenticationServiceI)(nil).OnStart))
}

// RemoveClient mocks base method.
func (m *MockAuthenticationServiceI) RemoveClient(arg0 models.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveClient", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveClient indicates an expected call of RemoveClient.
func (mr *MockAuthenticationServiceIMockRecorder) RemoveClient(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveClient", reflect.TypeOf((*MockAuthenticationServiceI)(nil).RemoveClient), arg0)
}

// RemoveEventListener mocks base method.
func (m *MockAuthenticationServiceI) RemoveEventListener(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveEventListener", arg0)
}

// RemoveEventListener indicates an expected call of RemoveEventListener.
func (mr *MockAuthenticationServiceIMockRecorder) RemoveEventListener(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEventListener", reflect.TypeOf((*MockAuthenticationServiceI)(nil).RemoveEventListener), arg0)
}

// SetParent mocks base method.
func (m *MockAuthenticationServiceI) SetParent(arg0 service.ServiceI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetParent", arg0)
}

// SetParent indicates an expected call of SetParent.
func (mr *MockAuthenticationServiceIMockRecorder) SetParent(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetParent", reflect.TypeOf((*MockAuthenticationServiceI)(nil).SetParent), arg0)
}

// Start mocks base method.
func (m *MockAuthenticationServiceI) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockAuthenticationServiceIMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Start))
}

// StripAuthFromMessage mocks base method.
func (m *MockAuthenticationServiceI) StripAuthFromMessage(arg0 *models.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StripAuthFromMessage", arg0)
}

// StripAuthFromMessage indicates an expected call of StripAuthFromMessage.
func (mr *MockAuthenticationServiceIMockRecorder) StripAuthFromMessage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StripAuthFromMessage", reflect.TypeOf((*MockAuthenticationServiceI)(nil).StripAuthFromMessage), arg0)
}

// UpgradeAuth mocks base method.
func (m *MockAuthenticationServiceI) UpgradeAuth(arg0 models.Key, arg1 models.RoleName, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradeAuth", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpgradeAuth indicates an expected call of UpgradeAuth.
func (mr *MockAuthenticationServiceIMockRecorder) UpgradeAuth(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeAuth", reflect.TypeOf((*MockAuthenticationServiceI)(nil).UpgradeAuth), arg0, arg1, arg2)
}

// ValidateAuthInMessage mocks base method.
func (m *MockAuthenticationServiceI) ValidateAuthInMessage(arg0 *models.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAuthInMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAuthInMessage indicates an expected call of ValidateAuthInMessage.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateAuthInMessage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAuthInMessage", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateAuthInMessage), arg0)
}

// ValidateClientForTopic mocks base method.
func (m *MockAuthenticationServiceI) ValidateClientForTopic(arg0 models.Key, arg1 models.MessageTopic) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateClientForTopic", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateClientForTopic indicates an expected call of ValidateClientForTopic.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateClientForTopic(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateClientForTopic", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateClientForTopic), arg0, arg1)
}

// ValidateSecret mocks base method.
func (m *MockAuthenticationServiceI) ValidateSecret(arg0 models.RoleName, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSecret", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateSecret indicates an expected call of ValidateSecret.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateSecret(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSecret", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateSecret), arg0, arg1)
}