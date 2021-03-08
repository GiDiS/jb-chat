package users

import (
	"jb_chat/pkg/events"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
)

type UsersListRequest struct {
	store.UserSearchCriteria
}

type UsersListResponse struct {
	events.ResultStatus
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
	events.ResultStatus
	User models.UserInfo `json:"user"`
}

func (l *UsersInfoResponse) SetUser(user models.UserInfo) {
	l.Ok = true
	l.User = user
}
