package main

import (
	"github.com/CameronHonis/chess-arbitrator/app"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	appService := app.BuildServices()
	appService.Start()
	wg.Wait()
}
