package server

import (
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

func StartWSServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			fmt.Println(connErr)
			return
		}
		_, addClientErr := GetUserClientsManager().AddNewClient(conn)
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
	GetLogManager().LogPrompt("server", prompt.SenderKey, prompt)
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
		GetLogManager().Log("server", fmt.Sprintf("unhandled prompt type %d", prompt.Type))
		return
	}
	if !parsedContent {
		GetLogManager().Log("server", fmt.Sprintf("could not parse prompt content for prompt type %d", prompt.Type))
		return
	}
}

func handleInitClientPrompt(clientKey string, content *InitClientPromptContent) {
	GetLogManager().Log("server", fmt.Sprintf("handling init client prompt for client %s", clientKey))
}

func handleSubscribeToTopicPrompt(clientKey string, content *SubscribeToTopicPromptContent) {
	alreadySubbedErr := GetUserClientsManager().SubscribeClientTo(clientKey, content.Topic)
	if alreadySubbedErr != nil {
		GetLogManager().Log("server", alreadySubbedErr)
	}
}

func handleUnsubscribeToTopicPrompt(clientKey string, content *UnsubscribeToTopicPromptContent) {
	err := GetUserClientsManager().UnsubClientFrom(clientKey, content.Topic)
	if err != nil {
		GetLogManager().Log("server", err)
	}
}

func handleTransferMessagePrompt(clientKey string, content *TransferMessagePromptContent) {
	subbedClientKeys := GetUserClientsManager().GetClientKeysSubscribedToTopic(content.Message.Topic)
	for _, subbedClientKey := range subbedClientKeys.Flatten() {
		if subbedClientKey != clientKey {
			subbedClient, err := GetUserClientsManager().GetClientFromKey(subbedClientKey)
			if err != nil {
				GetLogManager().Log("server", err)
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
