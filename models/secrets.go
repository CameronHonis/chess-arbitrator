package models

type Secret string

const (
	SECRET_ENV                    Secret = "ENV"
	SECRET_BOT_CLIENT_SECRET      Secret = "BOT_CLIENT_SECRET"
	SECRET_AUTH_KEY_MINS_TO_STALE Secret = "AUTH_KEY_MINS_TO_STALE"
)
