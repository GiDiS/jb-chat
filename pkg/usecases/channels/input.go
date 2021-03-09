package channels

import (
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
)

type ChannelsOneResult struct {
	events.ResultStatus
	Channel   *models.Channel  `json:"channel,omitempty"`
	ChannelId models.ChannelId `json:"cid"`
}

func (r *ChannelsOneResult) SetChannel(channel *models.Channel) {
	r.Ok = true
	r.Channel = channel
	if channel != nil {
		r.ChannelId = channel.Cid
	} else {
		r.ChannelId = models.NoChannel
	}
}

type ChannelsCreateRequest struct {
	Title   string       `json:"title"`
	Members []models.Uid `json:"members,omitempty"`
}

type ChannelsUpdateRequest struct {
	ChannelId models.ChannelId `json:"cid"`
	Title     models.ChannelId `json:"title"`
}

type ChannelsDeleteRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsJoinRequest struct {
	ChannelId models.ChannelId `json:"cid"`
	UserId    models.Uid       `json:"uid"`
}

type ChannelsLeaveRequest struct {
	ChannelId models.ChannelId `json:"cid"`
	UserId    models.Uid       `json:"uid"`
}

type ChannelsKickRequest struct {
	ChannelId models.ChannelId `json:"cid"`
	UserId    models.Uid       `json:"uid"`
}

type ChannelsGetListRequest struct {
	store.ChannelsSearchCriteria
}

type ChannelsListResponse struct {
	events.ResultStatus
	Channels []models.Channel `json:"channels"`
}

type ChannelsMembersRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsMembersResponse struct {
	events.ResultStatus
	Members []models.Uid `json:"members"`
}

type ChannelsGetInfoRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsGetDirectRequest struct {
	UserId   models.Uid `json:"uid"`
	Nickname string     `json:"nickname"`
}

func (l *ChannelsListResponse) SetChannels(channels []models.Channel) {
	l.Ok = true
	l.Channels = channels
}

func (l *ChannelsMembersResponse) SetMembers(members []models.Uid) {
	l.Ok = true
	l.Members = members
}
