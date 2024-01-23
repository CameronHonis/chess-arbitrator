package clients_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/msg_service"
	"github.com/CameronHonis/chess-arbitrator/sub_service"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"sync"
)

type ClientsManagerI interface {
	ServiceI
	GetClient(clientKey models.Key) (*models.Client, error)
	AddNewClient(conn *websocket.Conn) (*models.Client, error)
	AddClient(client *models.Client) error
	RemoveClient(client *models.Client) error
	BroadcastMessage(message *models.Message)
	DirectMessage(message *models.Message, clientKey models.Key) error
}

type ClientsManager struct {
	Service

	__dependencies__ Marker
	LogService       LoggerServiceI
	MsgService       msg_service.MessageServiceI
	SubService       sub_service.SubscriptionServiceI
	AuthService      auth.AuthenticationServiceI

	__state__         Marker
	clientByPublicKey map[models.Key]*models.Client
	mu                sync.Mutex
}

func NewClientsManager(config *ClientsManagerConfig) *ClientsManager {
	s := &ClientsManager{
		clientByPublicKey: make(map[models.Key]*models.Client),
	}
	s.Service = *NewService(s, config)

	return s
}

func (c *ClientsManager) OnBuild() {
	c.AddEventListener(CLIENT_CREATED, OnClientCreated)
	c.AddEventListener(auth.AUTH_UPGRADE_GRANTED, OnUpgradeAuthGranted)
	c.AddEventListener(matcher.CHALLENGE_REQUEST_FAILED, OnChallengeRequestFailed)
	c.AddEventListener(matcher.CHALLENGE_CREATED, OnChallengeCreated)
	c.AddEventListener(matcher.CHALLENGE_REVOKED, OnChallengeRevoked)
	c.AddEventListener(matcher.CHALLENGE_DENIED, OnChallengeDenied)
}

func (c *ClientsManager) AddNewClient(conn *websocket.Conn) (*models.Client, error) {
	client := auth.CreateClient(conn, c.CleanupClient)

	c.AuthService.AddClient(client.PublicKey())
	if err := c.AddClient(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ClientsManager) AddClient(client *models.Client) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.clientByPublicKey[client.PublicKey()]; ok {
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	c.clientByPublicKey[client.PublicKey()] = client
	go c.Dispatch(NewClientCreatedEvent(client))
	return nil
}

func (c *ClientsManager) RemoveClient(client *models.Client) error {
	pubKey := client.PublicKey()

	c.mu.Lock()
	if _, ok := c.clientByPublicKey[pubKey]; !ok {
		c.mu.Unlock()
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(c.clientByPublicKey, pubKey)
	c.mu.Unlock()

	c.SubService.UnsubClientFromAll(pubKey)
	return nil
}

func (c *ClientsManager) GetClient(clientKey models.Key) (*models.Client, error) {
	c.mu.Lock()
	client, ok := c.clientByPublicKey[clientKey]
	c.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (c *ClientsManager) BroadcastMessage(message *models.Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := c.SubService.ClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		client, err := c.GetClient(clientKey)
		if err != nil {
			c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), ALL_BUT_TEST_ENV)
			continue
		}
		writeErr := c.writeMessage(client, &msgCopy)
		if writeErr != nil {
			c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error broadcasting to client: %s", writeErr), ALL_BUT_TEST_ENV)
			continue
		}
	}
}

func (c *ClientsManager) DirectMessage(message *models.Message, clientKey models.Key) error {
	if message.Topic != "directMessage" && message.Topic != "" {
		return fmt.Errorf("direct messages expected to not have a topic, given %s", message.Topic)
	}
	msgCopy := *message
	msgCopy.Topic = "directMessage"
	client, clientErr := c.GetClient(clientKey)
	if clientErr != nil {
		return clientErr
	}
	return c.writeMessage(client, &msgCopy)
}

func (c *ClientsManager) CleanupClient(client *models.Client) {
	_ = c.AuthService.RemoveClient(client.PublicKey())
}

func (c *ClientsManager) listenForUserInput(client *models.Client) {
	if client.WSConn() == nil {
		c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("client %s did not establish a websocket connection", client.PublicKey()))
		return
	}
	for {
		_, rawMsg, readErr := client.WSConn().ReadMessage()
		_, clientErr := c.GetClient(client.PublicKey())
		if clientErr != nil {
			c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error listening on websocket: %s", clientErr), ALL_BUT_TEST_ENV)
			return
		}
		if readErr != nil {
			c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", readErr), ALL_BUT_TEST_ENV)
			// assume all readErrs are disconnects
			_ = c.RemoveClient(client)
			return
		}
		if err := c.readMessage(client.PublicKey(), rawMsg); err != nil {
			c.LogService.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", err), ALL_BUT_TEST_ENV)

		}
	}

}

func (c *ClientsManager) readMessage(clientKey models.Key, rawMsg []byte) error {
	c.LogService.Log(string(clientKey), ">> ", string(rawMsg))
	msg, unmarshalErr := models.UnmarshalToMessage(rawMsg)
	if unmarshalErr != nil {
		return fmt.Errorf("error unmarshalling message: %s", unmarshalErr)
	}
	if authErr := c.AuthService.ValidateAuthInMessage(msg); authErr != nil {
		return fmt.Errorf("error validating auth in message: %s", authErr)
	}
	c.AuthService.StripAuthFromMessage(msg)

	c.MsgService.HandleMessage(msg)
	c.BroadcastMessage(msg)
	return nil
}

func (c *ClientsManager) writeMessage(client *models.Client, msg *models.Message) error {
	msgJson, jsonErr := msg.Marshal()
	if jsonErr != nil {
		return jsonErr
	}
	c.LogService.Log(string(client.PublicKey()), "<< ", string(msgJson))
	return client.WSConn().WriteMessage(websocket.TextMessage, msgJson)
}
