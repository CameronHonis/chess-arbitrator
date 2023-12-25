package main

import (
	"github.com/CameronHonis/chess-arbitrator/server"
	"sync"
)

func main() {
	app := server.BuildServices()

	wg := sync.WaitGroup{}
	wg.Add(1)
	app.RouterService.StartWSServer()
	wg.Wait()
}
