package server

import "fmt"

func HandleMessage(msg *Message) {
	switch msg.Topic {
	case "findMatch":
		GetLogManager().LogGreen("handling find match message")
	case "move":
		GetLogManager().LogGreen("handling find match message")
	default:
		fmt.Println("unrecognized message topic: ", msg.Topic)
	}
	GetUserClientsManager().BroadcastMessage(msg)
}

//func handleFindMatchMessage
