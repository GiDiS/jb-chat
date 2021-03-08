package channels

import (
	"context"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
)

type Channels interface {
	GetList(ctx context.Context, request ChannelsGetListRequest) (resp ChannelsListResponse, err error)
	GetDirect(ctx context.Context, uid models.Uid, request ChannelsGetDirectRequest) (resp ChannelsOneResult, err error)
	GetMembers(ctx context.Context, uid models.Uid, request ChannelsMembersRequest) (resp ChannelsMembersResponse, err error)
	Get(ctx context.Context, cid models.ChannelId) (resp ChannelsOneResult, err error)
	Create(ctx context.Context, uid models.Uid, request ChannelsCreateRequest) (resp ChannelsOneResult, err error)
	Delete(ctx context.Context, uid models.Uid, request ChannelsDeleteRequest) (resp ChannelsOneResult, err error)
	Leave(ctx context.Context, uid models.Uid, request ChannelsLeaveRequest) (resp ChannelsOneResult, err error)
	Join(ctx context.Context, uid models.Uid, request ChannelsJoinRequest) (resp ChannelsOneResult, err error)
}

type channelsImpl struct {
	logger        logger.Logger
	channelsStore store.ChannelsStore
	membersStore  store.ChannelMembersStore
	usersStore    store.UsersStore
}

func NewChannels(
	logger logger.Logger,
	channelsStore store.ChannelsStore,
	membersStore store.ChannelMembersStore,
	usersStore store.UsersStore,
) *channelsImpl {
	return &channelsImpl{
		logger:        logger,
		channelsStore: channelsStore,
		membersStore:  membersStore,
		usersStore:    usersStore,
	}
}
