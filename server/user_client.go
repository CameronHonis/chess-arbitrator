package server

import (
	"github.com/gorilla/websocket"
)

type Key string

type Client struct {
	publicKey  Key
	privateKey Key
	conn       *websocket.Conn
	cleanup    func(*Client)
}

func NewClient(conn *websocket.Conn, cleanup func(*Client)) *Client {
	pubKey, priKey := GenerateKeyset()

	uc := &Client{
		publicKey:  pubKey,
		privateKey: priKey,
		conn:       conn,
		cleanup:    cleanup,
	}
	return uc
}

func (c *Client) PublicKey() Key {
	return c.publicKey
}

func (c *Client) WSConn() *websocket.Conn {
	return c.conn
}
