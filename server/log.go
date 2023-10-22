package server

import (
	"encoding/json"
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

func (lm *LogManager) Log(env string, msg interface{}) {
	lm.mu.Lock()
	fmt.Println("[", strings.ToUpper(env), "] ", msg)
	lm.mu.Unlock()
}

func (lm *LogManager) LogPrompt(env string, origin string, prompt *Prompt) {
	if LOG_INCOMING_PROMPTS {
		promptJson, err := json.Marshal(*prompt)
		lm.mu.Lock()
		if err != nil {
			fmt.Println("could not marshal json for ", *prompt, " while logging incoming prompt")
		} else {
			fmt.Println("[", strings.ToUpper(env), "] ", origin, " >> ", string(promptJson))
		}
		lm.mu.Unlock()
	}
}

func (lm *LogManager) LogMessage(env string, origin string, msg string) {
	if LOG_INCOMING_MESSAGES {
		lm.mu.Lock()
		fmt.Println("[", strings.ToUpper(env), "] ", origin, " >> ", msg)
		lm.mu.Unlock()
	}
}

func (lm *LogManager) LogClientActivity(env string, clientKey string, msg string) {
	lm.mu.Lock()
	fmt.Println("[", strings.ToUpper(env), "] ", clientKey, msg)
	lm.mu.Unlock()
}
