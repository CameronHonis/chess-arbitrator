package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	"github.com/gorilla/websocket"
	"strings"
)

type UserClient struct {
	publicKey  string
	privateKey string
	conn       *websocket.Conn
	cleanup    func(*UserClient)
}

func NewUserClient(conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := GenerateKeyset()

	uc := &UserClient{
		publicKey:  pubKey,
		privateKey: priKey,
		conn:       conn,
		cleanup:    cleanup,
	}
	return uc
}

func (uc *UserClient) PublicKey() string {
	return uc.publicKey
}

func (uc *UserClient) WSConn() *websocket.Conn {
	return uc.conn
}

func ClientKeyLogDecorator(clientKey string) string {
	concatKey := clientKey[:4] + ".." + clientKey[len(clientKey)-4:]
	return WrapCyan(fmt.Sprintf("%s", strings.ToUpper(concatKey)))
}
