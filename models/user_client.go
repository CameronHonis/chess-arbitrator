package models

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

func NewClient(pubKey, priKey Key, conn *websocket.Conn, cleanup func(*Client)) *Client {
	return &Client{
		pubKey,
		priKey,
		conn,
		cleanup,
	}
}

func (c *Client) PublicKey() Key {
	return c.publicKey
}

func (c *Client) PrivateKey() Key {
	return c.privateKey
}

func (c *Client) WSConn() *websocket.Conn {
	return c.conn
}
