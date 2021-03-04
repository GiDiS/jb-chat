package event

import (
	"jb_chat/pkg/models"
)

type ChannelsOneResult struct {
	ResultStatus
	Channel *models.Channel `json:"channel,omitempty"`
}

func (r *ChannelsOneResult) SetChannel(channel *models.Channel) {
	r.Ok = true
	r.Channel = channel
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
}

type ChannelsLeaveRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsKickRequest struct {
	ChannelId models.ChannelId `json:"cid"`
	UserId    models.ChannelId `json:"uid"`
}

type ChannelsGetListRequest struct {
}

type ChannelsListResponse struct {
	ResultStatus
	Channels []models.Channel `json:"channels"`
}

type ChannelsMembersRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsMembersResponse struct {
	ResultStatus
	Members []models.Uid `json:"members"`
}

type ChannelsGetInfoRequest struct {
	ChannelId models.ChannelId `json:"cid"`
}

type ChannelsGetDirectRequest struct {
	DirectUserId models.Uid `json:"uid"`
}

func (l *ChannelsListResponse) SetChannels(channels []models.Channel) {
	l.Ok = true
	l.Channels = channels
}

func (l *ChannelsMembersResponse) SetMembers(members []models.Uid) {
	l.Ok = true
	l.Members = members
}
