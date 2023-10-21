package server

import (
	"encoding/json"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

func StartWSServer() {
	userClientsManager, _ = NewUserClientsManager()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			fmt.Println(connErr)
			return
		}
		_, addClientErr := userClientsManager.AddNewClient(conn)
		if addClientErr != nil {
			fmt.Println(addClientErr)
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func upgradeToWSCon(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func handlePrompt(prompt *Prompt) {
	if LOG_INCOMING_PROMPTS {
		promptJson, err := json.Marshal(*prompt)
		userClientsManager.stdoutMutex.Lock()
		if err != nil {
			fmt.Println("could not marshal json for ", *prompt, " while logging incoming prompt")
		} else {
			fmt.Println("[SERVER]", prompt.SenderKey, " >> ", string(promptJson))
		}
		userClientsManager.stdoutMutex.Unlock()
	}
	parsedContent := false
	switch prompt.Type {
	case PROMPT_TYPE_INIT_CLIENT:
		if content, ok := prompt.Content.(InitClientPromptContent); ok {
			handleInitClientPrompt(prompt.SenderKey, &content)
			parsedContent = true
		}
	case PROMPT_TYPE_SUBSCRIBE_TO_TOPIC:
		if content, ok := prompt.Content.(SubscribeToTopicPromptContent); ok {
			handleSubscribeToTopicPrompt(prompt.SenderKey, &content)
			parsedContent = true
		}
	case PROMPT_TYPE_UNSUBSCRIBE_TO_TOPIC:
		if content, ok := prompt.Content.(UnsubscribeToTopicPromptContent); ok {
			handleUnsubscribeToTopicPrompt(prompt.SenderKey, &content)
			parsedContent = true
		}
	case PROMPT_TYPE_TRANSFER_MESSAGE:
		if content, ok := prompt.Content.(*TransferMessagePromptContent); ok {
			handleTransferMessagePrompt(prompt.SenderKey, content)
			parsedContent = true
		}
	default:
		userClientsManager.stdoutMutex.Lock()
		fmt.Printf("unhandled prompt type %d\n", prompt.Type)
		userClientsManager.stdoutMutex.Unlock()
		return
	}
	if !parsedContent {
		userClientsManager.stdoutMutex.Lock()
		fmt.Println("could not parse prompt content for prompt type ", prompt.Type)
		userClientsManager.stdoutMutex.Unlock()
		return
	}
}

func handleInitClientPrompt(clientKey string, content *InitClientPromptContent) {
	userClientsManager.stdoutMutex.Lock()
	fmt.Println("initialized client with key ", clientKey)
	userClientsManager.stdoutMutex.Unlock()
}

func handleSubscribeToTopicPrompt(clientKey string, content *SubscribeToTopicPromptContent) {
	alreadySubbedErr := userClientsManager.SubscribeClientTo(clientKey, content.Topic)
	if alreadySubbedErr != nil {
		userClientsManager.stdoutMutex.Lock()
		fmt.Println(alreadySubbedErr)
		userClientsManager.stdoutMutex.Unlock()
	}
}

func handleUnsubscribeToTopicPrompt(clientKey string, content *UnsubscribeToTopicPromptContent) {
	err := userClientsManager.UnsubClientFrom(clientKey, content.Topic)
	if err != nil {
		userClientsManager.stdoutMutex.Lock()
		fmt.Println(err)
		userClientsManager.stdoutMutex.Unlock()
	}
}

func handleTransferMessagePrompt(clientKey string, content *TransferMessagePromptContent) {
	subbedClientKeys := userClientsManager.GetClientKeysSubscribedToTopic(content.Message.Topic)
	for _, subbedClientKey := range subbedClientKeys.Flatten() {
		if subbedClientKey != clientKey {
			subbedClient, err := userClientsManager.GetClientFromKey(subbedClientKey)
			if err != nil {
				userClientsManager.stdoutMutex.Lock()
				fmt.Println(err)
				userClientsManager.stdoutMutex.Unlock()
				continue
			}
			subbedClient.InChannel() <- &Prompt{
				Type:      PROMPT_TYPE_TRANSFER_MESSAGE,
				SenderKey: clientKey,
				Content:   content,
			}
		}
	}
}
