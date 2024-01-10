package main

import (
	"github.com/CameronHonis/chess-arbitrator/app"
	"sync"
)

func main() {
	app := app.BuildServices()

	wg := sync.WaitGroup{}
	wg.Add(1)
	app.RouterService.StartWSServer()
	wg.Wait()
}
