package users

import "jb_chat/pkg/events"

const (
	UsersGetList events.Type = "users.get-list"
	UsersList    events.Type = "users.list"
	UsersGetInfo events.Type = "users.get-info"
	UsersInfo    events.Type = "users.info"
)
