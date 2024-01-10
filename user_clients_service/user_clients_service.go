package user_clients_service

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	"github.com/CameronHonis/chess-arbitrator/helpers"
	"github.com/CameronHonis/chess-arbitrator/message_service"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"sync"
)

type UserClientsServiceI interface {
	ServiceI
	GetClient(clientKey models.Key) (*models.Client, error)
	AddNewClient(conn *websocket.Conn) (*models.Client, error)
	AddClient(client *models.Client) error
	RemoveClient(client *models.Client) error
	BroadcastMessage(message *models.Message)
	DirectMessage(message *models.Message, clientKey models.Key) error
}

type UserClientsService struct {
	Service

	__dependencies__ Marker
	LogService       LoggerServiceI
	MsgService       message_service.MessageServiceI
	SubService       subscription_service.SubscriptionServiceI
	AuthService      auth_service.AuthenticationServiceI

	__state__         Marker
	mu                sync.Mutex
	clientByPublicKey map[models.Key]*models.Client
}

func NewUserClientsService(config *UserClientsServiceConfig) *UserClientsService {
	userClientsService := &UserClientsService{
		clientByPublicKey: make(map[models.Key]*models.Client),
	}
	userClientsService.Service = *NewService(userClientsService, config)
	return userClientsService
}

func (uc *UserClientsService) AddNewClient(conn *websocket.Conn) (*models.Client, error) {
	client := helpers.CreateClient(conn, uc.CleanupClient)

	if err := uc.AddClient(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (uc *UserClientsService) AddClient(client *models.Client) error {
	uc.mu.Lock()
	if _, ok := uc.clientByPublicKey[client.PublicKey()]; ok {
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	uc.clientByPublicKey[client.PublicKey()] = client
	uc.mu.Unlock()
	go uc.listenForUserInput(client)
	return nil
}

func (uc *UserClientsService) RemoveClient(client *models.Client) error {
	pubKey := client.PublicKey()

	uc.mu.Lock()
	if _, ok := uc.clientByPublicKey[pubKey]; !ok {
		uc.mu.Unlock()
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(uc.clientByPublicKey, pubKey)
	uc.mu.Unlock()

	uc.SubService.UnsubClientFromAll(pubKey)
	return nil
}

func (uc *UserClientsService) GetClient(clientKey models.Key) (*models.Client, error) {
	defer uc.mu.Unlock()
	uc.mu.Lock()
	client, ok := uc.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (uc *UserClientsService) BroadcastMessage(message *models.Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := uc.SubService.ClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		client, err := uc.GetClient(clientKey)
		if err != nil {
			uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), ALL_BUT_TEST_ENV)
			continue
		}
		writeErr := uc.writeMessage(client, &msgCopy)
		if writeErr != nil {
			uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error broadcasting to client: %s", writeErr), ALL_BUT_TEST_ENV)
			continue
		}
	}
}

func (uc *UserClientsService) DirectMessage(message *models.Message, clientKey models.Key) error {
	if message.Topic != "directMessage" && message.Topic != "" {
		return fmt.Errorf("direct messages expected to not have a topic, given %s", message.Topic)
	}
	msgCopy := *message
	msgCopy.Topic = "directMessage"
	client, clientErr := uc.GetClient(clientKey)
	if clientErr != nil {
		return clientErr
	}
	return uc.writeMessage(client, &msgCopy)
}

func (uc *UserClientsService) CleanupClient(client *models.Client) {
	_ = uc.AuthService.RemoveClient(client.PublicKey())
}

func (uc *UserClientsService) listenForUserInput(client *models.Client) {
	if client.WSConn() == nil {
		uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("client %s did not establish a websocket connection", client.PublicKey()))
		return
	}
	for {
		_, rawMsg, readErr := client.WSConn().ReadMessage()
		_, clientErr := uc.GetClient(client.PublicKey())
		if clientErr != nil {
			uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error listening on websocket: %s", clientErr), ALL_BUT_TEST_ENV)
			return
		}
		if readErr != nil {
			uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", readErr), ALL_BUT_TEST_ENV)
			// assume all readErrs are disconnects
			_ = uc.RemoveClient(client)
			return
		}
		if err := uc.readMessage(client.PublicKey(), rawMsg); err != nil {
			uc.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", err), ALL_BUT_TEST_ENV)

		}
	}

}

func (uc *UserClientsService) readMessage(clientKey models.Key, rawMsg []byte) error {
	uc.LogService.Log(string(clientKey), ">> ", string(rawMsg))
	msg, unmarshalErr := models.UnmarshalToMessage(rawMsg)
	if unmarshalErr != nil {
		return fmt.Errorf("error unmarshalling message: %s", unmarshalErr)
	}
	if authErr := uc.AuthService.ValidateAuthInMessage(msg); authErr != nil {
		return fmt.Errorf("error validating auth in message: %s", authErr)
	}
	uc.AuthService.StripAuthFromMessage(msg)

	uc.MsgService.HandleMessage(msg)
	uc.BroadcastMessage(msg)
	return nil
}

func (uc *UserClientsService) writeMessage(client *models.Client, msg *models.Message) error {
	client, err := uc.GetClient(msg.SenderKey)
	if err != nil {
		return err
	}
	msgJson, jsonErr := msg.Marshal()
	if jsonErr != nil {
		return jsonErr
	}
	uc.LogService.Log(string(client.PublicKey()), "<< ", string(msgJson))
	return client.WSConn().WriteMessage(websocket.TextMessage, msgJson)
}
