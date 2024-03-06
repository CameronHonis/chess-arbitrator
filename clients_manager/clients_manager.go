package clients_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	mm "github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/models"
	sub "github.com/CameronHonis/chess-arbitrator/sub_service"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"sync"
)

type ClientsManagerI interface {
	service.ServiceI

	AddConn(conn *websocket.Conn)
	BroadcastMessage(message *models.Message)
	DirectMessage(message *models.Message, clientKey models.Key) error
}

type ClientsManager struct {
	service.Service

	__dependencies__   marker.Marker
	Logger             log.LoggerServiceI
	SubService         sub.SubscriptionServiceI
	AuthService        auth.AuthenticationServiceI
	MatchmakingService mm.MatchmakingServiceI
	MatcherService     matcher.MatcherServiceI

	__state__    marker.Marker
	connByPubKey map[models.Key]*websocket.Conn
	mu           sync.Mutex
}

func NewClientsManager(config *ClientsManagerConfig) *ClientsManager {
	s := &ClientsManager{
		connByPubKey: make(map[models.Key]*websocket.Conn),
	}
	s.Service = *service.NewService(s, config)

	return s
}

func (c *ClientsManager) OnBuild() {
	c.AddEventListener(auth.CREDS_CHANGED, OnCredsChanged)
	c.AddEventListener(auth.CREDS_VETTED, OnCredsVetted)
	c.AddEventListener(auth.ROLE_SWITCHED, OnUpgradeAuthGranted)
	c.AddEventListener(matcher.CHALLENGE_REQUEST_FAILED, OnChallengeRequestFailed)
	c.AddEventListener(matcher.CHALLENGE_CREATED, OnChallengeCreated)
	c.AddEventListener(matcher.CHALLENGE_REVOKED, OnChallengeRevoked)
	c.AddEventListener(matcher.CHALLENGE_DENIED, OnChallengeDenied)
	c.AddEventListener(matcher.CHALLENGE_ACCEPTED, OnChallengeAccepted)
	c.AddEventListener(matcher.CHALLENGE_ACCEPT_FAILED, OnChallengeAcceptFailed)
	c.AddEventListener(matcher.MATCH_CREATED, OnMatchCreated)
	c.AddEventListener(matcher.MATCH_CREATION_FAILED, OnMatchCreationFailed)
	c.AddEventListener(matcher.MATCH_UPDATED, OnMatchUpdated)
	c.AddEventListener(matcher.MATCH_ENDED, OnMatchEnded)
	c.AddEventListener(matcher.MOVE_FAILURE, OnMoveFailed)
}

func (c *ClientsManager) AddConn(conn *websocket.Conn) {
	go c.listenOnConn(conn)
}

func (c *ClientsManager) BroadcastMessage(message *models.Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := c.SubService.ClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		conn, err := c.getConnByKey(clientKey)
		if err != nil {
			c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), log.ALL_BUT_TEST_ENV)
			continue
		}
		writeErr := c.writeMessage(clientKey, conn, &msgCopy)
		if writeErr != nil {
			c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error broadcasting to client: %s", writeErr), log.ALL_BUT_TEST_ENV)
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
	conn, err := c.getConnByKey(clientKey)
	if err != nil {
		return fmt.Errorf("unable to send DM: %s", err)
	}
	return c.writeMessage(clientKey, conn, &msgCopy)
}

func (c *ClientsManager) registerConn(pubKey models.Key, conn *websocket.Conn) error {
	if existingConn, _ := c.getConnByKey(pubKey); existingConn != nil {
		return fmt.Errorf("client already registered")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connByPubKey[pubKey] = conn
	return nil
}

func (c *ClientsManager) deregisterConn(pubKey models.Key) error {
	if _, err := c.getConnByKey(pubKey); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.connByPubKey, pubKey)
	return nil
}

func (c *ClientsManager) getConnByKey(pubKey models.Key) (*websocket.Conn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn, ok := c.connByPubKey[pubKey]
	if !ok {
		return nil, fmt.Errorf("no client with key %s", pubKey)
	}
	return conn, nil
}

func (c *ClientsManager) listenOnConn(conn *websocket.Conn) {
	var clientKey models.Key = "??"
	for {
		_, rawMsg, readErr := conn.ReadMessage()
		if readErr != nil {
			c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", readErr), log.ALL_BUT_TEST_ENV)
			if deregErr := c.deregisterConn(clientKey); deregErr != nil {
				c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error deregistering client: %s", deregErr), log.ALL_BUT_TEST_ENV)
			}
			return
		}
		c.Logger.Log(string(clientKey), ">> ", string(rawMsg))
		msg, unmarshalErr := models.UnmarshalToMessage(rawMsg)
		if unmarshalErr != nil {
			c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error unmarshalling message: %s", unmarshalErr))
			continue
		}

		// NOTE: this msg type is special since it requires the connection and doesn't require auth vetting
		if msg.ContentType == models.CONTENT_TYPE_REFRESH_AUTH {
			pubKey, refreshErr := HandleRefreshAuthMessage(c, msg, conn)
			if refreshErr != nil {
				c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("could not refresh creds: %s", refreshErr.Error()))
			} else {
				clientKey = pubKey
			}
			continue
		}

		if authErr := c.AuthService.VetAuthInMessage(msg); authErr != nil {
			sendDeps := NewSendDirectDeps(c.DirectMessage, clientKey)
			_ = SendInvalidAuth(sendDeps)
			c.Logger.LogRed(fmt.Sprintf("error validating auth in message: %s", authErr))
			continue
		} else if msg.SenderKey != "" {
			clientKey = msg.SenderKey
		}
		c.AuthService.StripAuthFromMessage(msg)

		if err := c.handleMsg(clientKey, msg); err != nil {
			c.Logger.LogRed(models.ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", err), log.ALL_BUT_TEST_ENV)
		}
	}
}

func (c *ClientsManager) handleMsg(clientKey models.Key, msg *models.Message) error {

	config := c.Config().(*ClientsManagerConfig)
	if msgHandler := config.HandlerByContentType(msg.ContentType); msgHandler != nil {
		if handlerErr := msgHandler(c, msg); handlerErr != nil {
			c.Logger.LogRed(models.ENV_CLIENT_MNGR, fmt.Sprintf("error handling msg \n\t%+v\n\n\t%s", *msg, handlerErr))
		}
	} else {
		return fmt.Errorf("no handler configured for msg %s", msg)
	}
	c.BroadcastMessage(msg)
	return nil
}

func (c *ClientsManager) writeMessage(pubkey models.Key, conn *websocket.Conn, msg *models.Message) error {
	msgJson, jsonErr := msg.Marshal()
	if jsonErr != nil {
		return jsonErr
	}
	c.Logger.Log(string(pubkey), "<< ", string(msgJson))

	c.mu.Lock()
	defer c.mu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, msgJson)
}
