package store

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"time"
)

type ChannelsStore interface {
	// CreateDirect create direct channel  between users
	CreateDirect(context.Context, models.Uid, models.Uid) (models.ChannelId, error)
	// GetDirect get direct channel  between users
	GetDirect(context.Context, models.Uid, models.Uid) (models.ChannelId, error)
	// CreatePublic create public chat room
	CreatePublic(context.Context, models.Uid, string) (models.ChannelId, error)
	// Delete complete delete channel
	Delete(context.Context, models.ChannelId) error
	// Get get channel info
	Get(context.Context, models.ChannelId) (models.Channel, error)
	// Find find channels by criteria
	Find(context.Context, ChannelsSearchCriteria) ([]models.Channel, error)
	// Estimate channels count by criteria
	Estimate(context.Context, ChannelsSearchCriteria) (uint64, error)
	// SetLastMessage set last message id in channel
	SetLastMessage(context.Context, models.ChannelId, models.MessageId, time.Time) error
	// IncMembersCount increase members count in channel
	IncMembersCount(context.Context, models.ChannelId, int) error
}

type ChannelMembersStore interface {
	// MemberOf check what user is member of channel
	MemberOf(context.Context, models.Uid) ([]models.ChannelId, error)
	// Members list users in channel
	Members(context.Context, models.ChannelId) ([]models.Uid, error)
	// Join add user to channel
	Join(context.Context, models.ChannelId, models.Uid) error
	// Leave remove user from channel
	Leave(context.Context, models.ChannelId, models.Uid) error
	// SetLastSeen set last seen by user message in the channel
	SetLastSeen(context.Context, models.ChannelId, models.Uid, models.MessageId) error
	// GetLastSeen get last seen by user message in the channel
	GetLastSeen(ctx context.Context, cid models.ChannelId, uid models.Uid) (msgId models.MessageId, err error)
}

type ChannelsSearchCriteria struct {
	ChannelIds []models.ChannelId `json:"channel_ids,omitempty"`
	HasMember  models.Uid         `json:"has_members,omitempty"`
	Title      string             `json:"title,omitempty"`
	Type       models.ChannelType `json:"type,omitempty"`
	Limits     models.Limits      `json:"limits,omitempty"`
}
