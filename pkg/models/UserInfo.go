package models

type UserStatus string

const (
	UserStatusUnknown UserStatus = "unknown"
	UserStatusOffline UserStatus = "offline"
	UserStatusOnline  UserStatus = "online"
	UserStatusAway    UserStatus = "away"
)

type UserInfo struct {
	User
	Status UserStatus `json:"status"`
}
