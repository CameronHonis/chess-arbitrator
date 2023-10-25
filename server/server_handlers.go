package server

import "fmt"

func HandleMessage(msg *Message, clientKey string) {
	var handleMsgErr error
	switch msg.ContentType {
	case CONTENT_TYPE_FIND_MATCH:
		handleMsgErr = HandleFindMatchMessage(clientKey)
	}
	if handleMsgErr != nil {
		GetLogManager().LogRed("server", fmt.Sprintf("could not handle message \n\t%s\n\t%s", msg, handleMsgErr))
	}
	GetUserClientsManager().BroadcastMessage(msg)
}

func HandleFindMatchMessage(clientKey string) error {
	// TODO: query for elo, winStreak, lossStreak
	addClientErr := GetMatchmakingManager().AddClient(&ClientProfile{
		ClientKey:  clientKey,
		Elo:        1000,
		WinStreak:  0,
		LossStreak: 0,
	})
	if addClientErr != nil {
		return fmt.Errorf("could not add client %s to matchmaking pool: %s", clientKey, addClientErr)
	}
	return nil
}
