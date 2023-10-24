package server

type MessageTopic uint8

const LOG_INCOMING_PROMPTS = true
const LOG_INCOMING_MESSAGES = true

const (
	MESSAGE_TOPIC_NONE MessageTopic = iota
	MESSAGE_TOPIC_AUTH
	MESSAGE_TOPIC_INIT_BOT_MATCH
)

type PromptType uint8

const (
	PROMPT_TYPE_NONE PromptType = iota
	PROMPT_TYPE_INIT_CLIENT
	PROMPT_TYPE_SUBSCRIBE_TO_TOPIC
	PROMPT_TYPE_UNSUBSCRIBE_TO_TOPIC
	PROMPT_TYPE_TRANSFER_MESSAGE
)

type BotType uint8

const (
	BOT_TYPE_NONE BotType = iota
	BOT_TYPE_RANDOM
)