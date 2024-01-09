package main

import (
	"sync"
)

func main() {
	app := BuildServices()

	wg := sync.WaitGroup{}
	wg.Add(1)
	app.RouterService.StartWSServer()
	wg.Wait()
}
