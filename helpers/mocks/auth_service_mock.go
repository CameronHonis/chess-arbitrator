// Code generated by MockGen. DO NOT EDIT.
// Source: ../auth/auth_service.go
//
// Generated by this command:
//
//	mockgen -source=../auth/auth_service.go -destination mocks/auth_service_mock.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

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
func (m *MockAuthenticationServiceI) AddClient(clientKey models.Key) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClient", clientKey)
}

// AddClient indicates an expected call of AddClient.
func (mr *MockAuthenticationServiceIMockRecorder) AddClient(clientKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClient", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddClient), clientKey)
}

// AddDependency mocks base method.
func (m *MockAuthenticationServiceI) AddDependency(service service.ServiceI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDependency", service)
}

// AddDependency indicates an expected call of AddDependency.
func (mr *MockAuthenticationServiceIMockRecorder) AddDependency(service any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDependency", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddDependency), service)
}

// AddEventListener mocks base method.
func (m *MockAuthenticationServiceI) AddEventListener(eventVariant service.EventVariant, fn service.EventHandler) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEventListener", eventVariant, fn)
	ret0, _ := ret[0].(int)
	return ret0
}

// AddEventListener indicates an expected call of AddEventListener.
func (mr *MockAuthenticationServiceIMockRecorder) AddEventListener(eventVariant, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventListener", reflect.TypeOf((*MockAuthenticationServiceI)(nil).AddEventListener), eventVariant, fn)
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

// Build mocks base method.
func (m *MockAuthenticationServiceI) Build() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Build")
}

// Build indicates an expected call of Build.
func (mr *MockAuthenticationServiceIMockRecorder) Build() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Build))
}

// ClientKeysByRole mocks base method.
func (m *MockAuthenticationServiceI) ClientKeysByRole(roleName models.RoleName) *set.Set[models.Key] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientKeysByRole", roleName)
	ret0, _ := ret[0].(*set.Set[models.Key])
	return ret0
}

// ClientKeysByRole indicates an expected call of ClientKeysByRole.
func (mr *MockAuthenticationServiceIMockRecorder) ClientKeysByRole(roleName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientKeysByRole", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ClientKeysByRole), roleName)
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
func (m *MockAuthenticationServiceI) Dispatch(event service.EventI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", event)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockAuthenticationServiceIMockRecorder) Dispatch(event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockAuthenticationServiceI)(nil).Dispatch), event)
}

// GetRole mocks base method.
func (m *MockAuthenticationServiceI) GetRole(clientKey models.Key) (models.RoleName, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", clientKey)
	ret0, _ := ret[0].(models.RoleName)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockAuthenticationServiceIMockRecorder) GetRole(clientKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockAuthenticationServiceI)(nil).GetRole), clientKey)
}

// OnBuild mocks base method.
func (m *MockAuthenticationServiceI) OnBuild() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnBuild")
}

// OnBuild indicates an expected call of OnBuild.
func (mr *MockAuthenticationServiceIMockRecorder) OnBuild() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnBuild", reflect.TypeOf((*MockAuthenticationServiceI)(nil).OnBuild))
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
func (m *MockAuthenticationServiceI) RemoveClient(clientKey models.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveClient", clientKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveClient indicates an expected call of RemoveClient.
func (mr *MockAuthenticationServiceIMockRecorder) RemoveClient(clientKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveClient", reflect.TypeOf((*MockAuthenticationServiceI)(nil).RemoveClient), clientKey)
}

// RemoveEventListener mocks base method.
func (m *MockAuthenticationServiceI) RemoveEventListener(eventId int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveEventListener", eventId)
}

// RemoveEventListener indicates an expected call of RemoveEventListener.
func (mr *MockAuthenticationServiceIMockRecorder) RemoveEventListener(eventId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEventListener", reflect.TypeOf((*MockAuthenticationServiceI)(nil).RemoveEventListener), eventId)
}

// SetParent mocks base method.
func (m *MockAuthenticationServiceI) SetParent(parent service.ServiceI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetParent", parent)
}

// SetParent indicates an expected call of SetParent.
func (mr *MockAuthenticationServiceIMockRecorder) SetParent(parent any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetParent", reflect.TypeOf((*MockAuthenticationServiceI)(nil).SetParent), parent)
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
func (m *MockAuthenticationServiceI) StripAuthFromMessage(msg *models.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StripAuthFromMessage", msg)
}

// StripAuthFromMessage indicates an expected call of StripAuthFromMessage.
func (mr *MockAuthenticationServiceIMockRecorder) StripAuthFromMessage(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StripAuthFromMessage", reflect.TypeOf((*MockAuthenticationServiceI)(nil).StripAuthFromMessage), msg)
}

// UpgradeAuth mocks base method.
func (m *MockAuthenticationServiceI) UpgradeAuth(clientKey models.Key, roleName models.RoleName, secret string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradeAuth", clientKey, roleName, secret)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpgradeAuth indicates an expected call of UpgradeAuth.
func (mr *MockAuthenticationServiceIMockRecorder) UpgradeAuth(clientKey, roleName, secret any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeAuth", reflect.TypeOf((*MockAuthenticationServiceI)(nil).UpgradeAuth), clientKey, roleName, secret)
}

// ValidateAuthInMessage mocks base method.
func (m *MockAuthenticationServiceI) ValidateAuthInMessage(msg *models.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAuthInMessage", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAuthInMessage indicates an expected call of ValidateAuthInMessage.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateAuthInMessage(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAuthInMessage", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateAuthInMessage), msg)
}

// ValidateClientForTopic mocks base method.
func (m *MockAuthenticationServiceI) ValidateClientForTopic(clientKey models.Key, topic models.MessageTopic) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateClientForTopic", clientKey, topic)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateClientForTopic indicates an expected call of ValidateClientForTopic.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateClientForTopic(clientKey, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateClientForTopic", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateClientForTopic), clientKey, topic)
}

// ValidateSecret mocks base method.
func (m *MockAuthenticationServiceI) ValidateSecret(roleName models.RoleName, secret string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSecret", roleName, secret)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateSecret indicates an expected call of ValidateSecret.
func (mr *MockAuthenticationServiceIMockRecorder) ValidateSecret(roleName, secret any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSecret", reflect.TypeOf((*MockAuthenticationServiceI)(nil).ValidateSecret), roleName, secret)
}
