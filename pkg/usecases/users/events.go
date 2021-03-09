package users

import "github.com/GiDiS/jb-chat/pkg/events"

const (
	UsersGetList events.Type = "users.get-list"
	UsersList    events.Type = "users.list"
	UsersGetInfo events.Type = "users.get-info"
	UsersInfo    events.Type = "users.info"
)
