package messages

import (
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/models"
)

type MessageRef struct {
	ChannelId models.ChannelId `json:"cid"`
	MessageId models.MessageId `json:"mid"`
}

type MessageOneResult struct {
	events.ResultStatus
	Msg  *models.Message  `json:"msg,omitempty"`
	User *models.UserInfo `json:"user,omitempty"`
}

type MessageCreateRequest struct {
	ChannelId models.ChannelId  `json:"cid"`
	UserId    models.Uid        `json:"uid"`
	ParentId  *models.MessageId `json:"pid,omitempty"`
	Body      string            `json:"body"`
}

type MessageUpdateRequest struct {
	MessageRef
	Body string `json:"body"`
}

type MessageListRequest struct {
	ChannelId models.ChannelId  `json:"cid"`
	ParentId  *models.MessageId `json:"pid"`
	After     *models.MessageId `json:"after"`
	Before    *models.MessageId `json:"before"`
	Limit     int
}

type MessageListResponse struct {
	events.ResultStatus
	Messages []models.Message               `json:"messages,omitempty"`
	Users    map[models.Uid]models.UserInfo `json:"users,omitempty"`
}

func (r *MessageListResponse) SetResult(msgs []models.Message, users []models.UserInfo) {
	r.Ok = true
	r.Messages = msgs
	r.Users = make(map[models.Uid]models.UserInfo)

	// dedup users
	for _, user := range users {
		if _, ok := r.Users[user.UserId]; !ok {
			r.Users[user.UserId] = user
		}
	}
}

func (r *MessageOneResult) SetMsg(msg *models.Message) {
	r.Msg = msg
	r.Ok = true
}
