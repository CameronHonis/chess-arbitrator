package main

import (
	"github.com/CameronHonis/chess-arbitrator/server"
)

func main() {
	server.ConfigLogger()
	server.StartWSServer()
}
