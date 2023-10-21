package server

type Prompt struct {
	Type      PromptType
	SenderKey string
	Content   interface{}
}

type InitClientPromptContent struct {
}

type SubscribeToTopicPromptContent struct {
	Topic MessageTopic
}

type UnsubscribeToTopicPromptContent struct {
	Topic MessageTopic
}

type TransferMessagePromptContent struct {
	Message *Message
}
