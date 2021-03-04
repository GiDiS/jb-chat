package event

import (
	"jb_chat/pkg/models"
)

type UsersListRequest struct {
}

type UsersListResponse struct {
	ResultStatus
	Users []models.UserInfo `json:"users"`
}

func (l *UsersListResponse) SetUsers(users []models.UserInfo) {
	l.Ok = true
	l.Users = users
}

type UsersInfoRequest struct {
	UserId models.Uid `json:"uid"`
}

type UsersInfoResponse struct {
	ResultStatus
	User models.UserInfo `json:"user"`
}

func (l *UsersInfoResponse) SetUser(user models.UserInfo) {
	l.Ok = true
	l.User = user
}
