package server

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Log struct {
	Env          string
	Msg          string
	ColorWrapper func(string) string
	Options      []LogOption
}

func (l *Log) String() string {
	if l.ColorWrapper == nil {
		return fmt.Sprintf("[%s] %s", strings.ToUpper(l.Env), l.Msg)
	} else {
		return l.ColorWrapper(fmt.Sprintf("[%s] %s", strings.ToUpper(l.Env), l.Msg))
	}
}

type LogBuilder struct {
	log *Log
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{
		log: &Log{},
	}
}

func (lb *LogBuilder) Env(env string) *LogBuilder {
	lb.log.Env = env
	return lb
}

func (lb *LogBuilder) Msg(msg string) *LogBuilder {
	lb.log.Msg = msg
	return lb
}

func (lb *LogBuilder) Msgs(msgs ...interface{}) *LogBuilder {
	var filteredMsg strings.Builder
	logOptions := make([]LogOption, 0, len(msgs))
	for _, msg := range msgs {
		if IsLogOption(msg) {
			logOptions = append(logOptions, LogOption(fmt.Sprintf("%v", msg)))
		} else {
			filteredMsg.WriteString(fmt.Sprintf("%v", msg))
		}
	}
	lb.log.Msg = filteredMsg.String()
	lb.log.Options = logOptions
	return lb
}

func (lb *LogBuilder) ColorWrapper(colorWrapper func(string) string) *LogBuilder {
	lb.log.ColorWrapper = colorWrapper
	return lb
}

func (lb *LogBuilder) Options(options ...LogOption) *LogBuilder {
	lb.log.Options = options
	return lb
}

func (lb *LogBuilder) Build() *Log {
	return lb.log
}

var logManager *LogManager

type LogManager struct {
	mu *sync.Mutex
}

func GetLogManager() *LogManager {
	if logManager != nil {
		return logManager
	}
	return &LogManager{
		mu: &sync.Mutex{},
	}
}

func (lm *LogManager) logLogWithLock(log *Log) {
	lm.mu.Lock()
	fmt.Println(log.String())
	lm.mu.Unlock()
}

func (lm *LogManager) canPrintInEnv(msgs ...interface{}) bool {
	for _, msg := range msgs {
		switch msg {
		case ONLY_TEST_ENV:
			return GetEnvName() == "test"
		case ONLY_PROD_ENV:
			return GetEnvName() == "prod"
		case ALL_BUT_TEST_ENV:
			return GetEnvName() != "test"
		}
	}
	return true
}

func (lm *LogManager) Log(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}

	log := NewLogBuilder().Env(env).Msgs(msgs...).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogRed(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(lm.WrapRed).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogGreen(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(lm.WrapGreen).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogMessage(env string, origin string, msg string) {
	if LOG_INCOMING_MESSAGES {
		lm.mu.Lock()
		fmt.Println(fmt.Sprintf("[%s] %s >> %s", strings.ToUpper(env), origin, msg))
		lm.mu.Unlock()
	}
}

func (lm *LogManager) LogClientActivity(env string, clientKey string, msg string) {
	lm.mu.Lock()
	fmt.Println("[", strings.ToUpper(env), "] ", clientKey, msg)
	lm.mu.Unlock()
}

func (lm *LogManager) WrapGreen(msg string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
}

func (lm *LogManager) WrapRed(msg string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
}

type LogOption string

const (
	ONLY_TEST_ENV    LogOption = "ONLY_TEST_ENV"
	ONLY_PROD_ENV              = "ONLY_PROD_ENV"
	ALL_BUT_TEST_ENV           = "ALL_BUT_TEST_ENV"
)

func IsLogOption(m interface{}) bool {
	return m == ONLY_TEST_ENV || m == ONLY_PROD_ENV || m == ALL_BUT_TEST_ENV
}

func GetEnvName() string {
	envName, hasEnvVar := os.LookupEnv("ENV")
	if !hasEnvVar {
		envName = "test"
	}
	return strings.ToLower(envName)
}
