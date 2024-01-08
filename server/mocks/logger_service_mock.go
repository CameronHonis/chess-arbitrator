package mocks

import (
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/stub"
)

type LoggerServiceMock struct {
	Stubbed[LoggerService]
	ServiceMock
}

func NewLoggerServiceMock(loggerService *LoggerService) *LoggerServiceMock {
	s := &LoggerServiceMock{}
	s.Stubbed = *NewStubbed(s, loggerService)
	s.ServiceMock = *NewServiceMock(&loggerService.Service)
	return s
}

func (s *LoggerServiceMock) Log(env string, msgs ...interface{}) {
	_ = s.Call("Log", env, msgs)
}
func (s *LoggerServiceMock) LogRed(env string, msgs ...interface{}) {
	_ = s.Call("LogRed", env, msgs)
}
func (s *LoggerServiceMock) LogGreen(env string, msgs ...interface{}) {
	_ = s.Call("LogGreen", env, msgs)
}
func (s *LoggerServiceMock) LogBlue(env string, msgs ...interface{}) {
	_ = s.Call("LogBlue", env, msgs)
}
func (s *LoggerServiceMock) LogYellow(env string, msgs ...interface{}) {
	_ = s.Call("LogYellow", env, msgs)
}
func (s *LoggerServiceMock) LogMagenta(env string, msgs ...interface{}) {
	_ = s.Call("LogMagenta", env, msgs)
}
func (s *LoggerServiceMock) LogCyan(env string, msgs ...interface{}) {
	_ = s.Call("LogCyan", env, msgs)
}
func (s *LoggerServiceMock) LogOrange(env string, msgs ...interface{}) {
	_ = s.Call("LogOrange", env, msgs)
}
func (s *LoggerServiceMock) LogBrown(env string, msgs ...interface{}) {
	_ = s.Call("LogBrown", env, msgs)
}
