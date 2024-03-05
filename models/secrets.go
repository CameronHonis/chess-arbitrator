package models

type Secret string

const (
	SECRET_ENV                    Secret = "SECRET_ENV"
	SECRET_BOT_CLIENT_SECRET      Secret = "SECRET_BOT_CLIENT_SECRET"
	SECRET_AUTH_KEY_MINS_TO_STALE Secret = "SECRET_AUTH_KEY_MINS_TO_STALE"
)
