package server

type MessageTopic uint8

const (
	AUTH MessageTopic = iota
	MATCHMAKING
)
