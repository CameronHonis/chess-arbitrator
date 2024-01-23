package main

import (
	"github.com/CameronHonis/chess-arbitrator/app"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	app := app.BuildServices()
	app.Build()
	app.Start()
	wg.Wait()
}
