package models

type Secret string

const (
	ENV                    Secret = "ENV"
	BOT_CLIENT_SECRET      Secret = "BOT_CLIENT_SECRET"
	AUTH_KEY_MINS_TO_STALE Secret = "AUTH_KEY_MINS_TO_STALE"
)
