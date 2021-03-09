package store

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

type ChannelsStore interface {
	CreateDirect(context.Context, models.Uid, models.Uid) (models.ChannelId, error)
	GetDirect(context.Context, models.Uid, models.Uid) (models.ChannelId, error)
	CreatePublic(context.Context, models.Uid, string) (models.ChannelId, error)
	Delete(context.Context, models.ChannelId) error
	Get(context.Context, models.ChannelId) (models.Channel, error)
	Find(context.Context, ChannelsSearchCriteria) ([]models.Channel, error)
	Estimate(context.Context, ChannelsSearchCriteria) (uint64, error)
}

type ChannelMembersStore interface {
	MemberOf(context.Context, models.Uid) ([]models.ChannelId, error)
	Members(context.Context, models.ChannelId) ([]models.Uid, error)
	Join(context.Context, models.ChannelId, models.Uid) error
	Leave(context.Context, models.ChannelId, models.Uid) error
}

type ChannelsSearchCriteria struct {
	ChannelIds []models.ChannelId `json:"channel_ids,omitempty"`
	HasMember  models.Uid         `json:"has_members,omitempty"`
	Title      string             `json:"title,omitempty"`
	Type       models.ChannelType `json:"type,omitempty"`
	Limits     models.Limits      `json:"limits,omitempty"`
}
