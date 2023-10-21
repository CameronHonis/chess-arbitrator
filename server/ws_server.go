package server

import (
	"encoding/json"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

func listenOnWsCon(con *websocket.Conn) {
	msg := AuthMessageContent{
		PublicKey:  "asdf",
		PrivateKey: "some-private-key",
	}
	msgJsonBytes, jsonParseErr := msg.toJsonBytes()
	if jsonParseErr != nil {
		fmt.Println(jsonParseErr)
		return
	}
	writeBytesErr := con.WriteMessage(websocket.TextMessage, msgJsonBytes)
	if writeBytesErr != nil {
		fmt.Println(writeBytesErr)
		return
	}
	for {
		messageType, p, err := con.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(messageType, string(p))
		if err := con.WriteMessage(websocket.TextMessage, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

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

	go func() {
		err := http.ListenAndServe(":8080", nil)
		fmt.Println("asdf")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
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

func handlePrompt(prompt *Prompt) error {
	if LOG_PROMPTS {
		promptJson, err := json.Marshal(*prompt)
		userClientsManager.stdoutMutex.Lock()
		if err != nil {
			fmt.Println("could not marshal json for ", *prompt)
		} else {
			fmt.Println(prompt.SenderKey, " >> ", string(promptJson))
		}
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
		if content, ok := prompt.Content.(TransferMessagePromptContent); ok {
			handleTransferMessagePrompt(prompt.SenderKey, &content)
			parsedContent = true
		}
	default:
		return fmt.Errorf("unhandled prompt type %d", prompt.Type)
	}
	if !parsedContent {
		return fmt.Errorf("could not parse prompt content for prompt type %d", prompt.Type)
	}
	return nil
}

func handleInitClientPrompt(clientKey string, content *InitClientPromptContent) {
	if LOG_PROMPTS {
		userClientsManager.stdoutMutex.Lock()
		fmt.Println("initialized client with key ", clientKey)
		userClientsManager.stdoutMutex.Unlock()
	}
}

func handleSubscribeToTopicPrompt(clientKey string, content *SubscribeToTopicPromptContent) {

}

func handleUnsubscribeToTopicPrompt(clientKey string, content *UnsubscribeToTopicPromptContent) {

}

func handleTransferMessagePrompt(clientKey string, content *TransferMessagePromptContent) {

}
