package models

type UserStatus uint8

const (
	UserStatusUnknown UserStatus = iota
	UserStatusOffline
	UserStatusOnline
	UserStatusAway
)

type UserInfo struct {
	User
	Status UserStatus `json:"status"`
}
