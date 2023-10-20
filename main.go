package main

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	"time"
)

func main() {
	server.StartWSServer()
	testSelectChannel()
}

func testSelectChannel() {
	channels := make([]chan int, 0)
	go func() {
		channels = append(channels, make(chan int))
		channels[0] <- 0
		time.Sleep(time.Second)
		channels = append(channels, make(chan int))
		channels = append(channels, make(chan int))
		channels[1] <- 1
		time.Sleep(time.Second)
		channels[2] <- 2
	}()
	for {
		for _, ch := range channels {
			select {
			case val := <-ch:
				fmt.Println(val)
			default:
				continue
			}
		}
	}
}
