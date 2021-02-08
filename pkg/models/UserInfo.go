package models

type UserStatus uint8

const (
	UserStatusUnknown UserStatus = iota
	UserStatusOffline
	UserStatusOnline
	UserStatusAway
)

type UserInfo struct {
	UserId    Uid        `json:"uid"`
	Status    UserStatus `json:"status"`
	Login     string     `json:"login"`
	ShownName string     `json:"shown_name"`
	Email     string     `json:"email"`
	Picture   string     `json:"picture,omitempty"`
}
