package main

import (
	"github.com/CameronHonis/chess-arbitrator/app_service"
	"sync"
)

func main() {
	app := app_service.BuildServices()

	wg := sync.WaitGroup{}
	wg.Add(1)
	app.RouterService.StartWSServer()
	wg.Wait()
}
