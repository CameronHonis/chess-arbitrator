package models

type RoleName string

const (
	PLEB RoleName = "PLEB"
	BOT           = "BOT"
)

var ENV_NAME_BY_ROLE_NAME = map[RoleName]string{
	PLEB: "PLEB_SECRET",
	BOT:  "BOT_CLIENT_SECRET",
}
