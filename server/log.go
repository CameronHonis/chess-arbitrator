package server

import (
	"fmt"
	"strings"
	"sync"
)

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

func (lm *LogManager) formatMessage(env string, msg ...interface{}) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] ", strings.ToUpper(env)))
	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%v", m))
	}
	return sb.String()
}

func (lm *LogManager) logWithLock(logMsg string) {
	lm.mu.Lock()
	fmt.Println(logMsg)
	lm.mu.Unlock()
}

func (lm *LogManager) Log(env string, msg ...interface{}) {
	lm.logWithLock(lm.formatMessage(env, msg...))
}

func (lm *LogManager) LogRed(env string, msg ...interface{}) {
	coloredString := fmt.Sprintf("\x1b[31m%s\x1b[0m", lm.formatMessage(env, msg...))
	lm.logWithLock(coloredString)
}

func (lm *LogManager) LogGreen(env string, msg ...interface{}) {
	coloredString := fmt.Sprintf("\x1b[32m%s\x1b[0m", lm.formatMessage(env, msg...))
	lm.logWithLock(coloredString)
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
