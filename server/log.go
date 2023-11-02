package server

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

type Log struct {
	Env          string
	Msg          string
	ColorWrapper func(string) string
	Options      []LogOption
}

func (l *Log) formatEnv() string {
	upperEnv := strings.ToUpper(l.Env)
	if upperEnv == "SERVER" {
		return WrapGreen(fmt.Sprintf("[%s]", upperEnv))
	} else if upperEnv == "CLIENT" {
		return WrapBlue(fmt.Sprintf("[%s]", upperEnv))
	} else if upperEnv == "MATCHMAKING" {
		return WrapMagenta(fmt.Sprintf("[%s]", upperEnv))
	} else if utf8.RuneCountInString(upperEnv) == 64 {
		envChars := []rune(l.Env)
		shortenedEnv := fmt.Sprintf("[%s..%s]", string(envChars[:4]), string(envChars[60:]))
		return WrapCyan(shortenedEnv)
	}
	return fmt.Sprintf("[%s]", upperEnv)
}

func (l *Log) String() string {
	if l.ColorWrapper == nil {
		return fmt.Sprintf("%s %s", l.formatEnv(), l.Msg)
	} else {
		return fmt.Sprintf("%s %s", l.formatEnv(), l.ColorWrapper(l.Msg))
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
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapRed).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogGreen(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapGreen).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogBlue(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapBlue).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogYellow(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapYellow).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogMagenta(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapMagenta).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogCyan(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapCyan).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogOrange(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapOrange).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogBrown(env string, msgs ...interface{}) {
	if !lm.canPrintInEnv(msgs...) {
		return
	}
	log := NewLogBuilder().Env(env).Msgs(msgs...).ColorWrapper(WrapBrown).Build()
	lm.logLogWithLock(log)
}

func (lm *LogManager) LogMessage(remote string, isIncoming bool, msg string) {
	if isIncoming {
		if LOG_INCOMING_MESSAGES {
			log := NewLogBuilder().Msg(fmt.Sprintf(">> %s", msg)).Env(remote).Build()
			lm.logLogWithLock(log)
		}
	} else {
		if LOG_OUTGOING_MESSAGES {
			log := NewLogBuilder().Msg(fmt.Sprintf("<< %s", msg)).Env(remote).Build()
			lm.logLogWithLock(log)
		}
	}
}

func WrapGreen(msg string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
}

func WrapRed(msg string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
}

func WrapYellow(msg string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg)
}

func WrapBlue(msg string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", msg)
}

func WrapMagenta(msg string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", msg)
}

func WrapCyan(msg string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", msg)
}

func WrapOrange(msg string) string {
	return fmt.Sprintf("\x1b[38;5;208m%s\x1b[0m", msg)
}

func WrapBrown(msg string) string {
	return fmt.Sprintf("\x1b[38;5;130m%s\x1b[0m", msg)
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
