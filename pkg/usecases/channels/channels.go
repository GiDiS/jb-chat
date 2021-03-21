package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
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
	GetLastSeen(ctx context.Context, cid models.ChannelId, uid models.Uid) (models.MessageId, error)
	SetLastSeen(ctx context.Context, cid models.ChannelId, uid models.Uid, mid models.MessageId) error
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
