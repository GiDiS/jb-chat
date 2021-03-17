package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/usecases"
)

func (c *channelsImpl) GetDirect(ctx context.Context, uid models.Uid, request ChannelsGetDirectRequest) (resp ChannelsOneResult, err error) {

	filter := store.UserSearchCriteria{}
	if request.UserId > 0 {
		filter.Uids = []models.Uid{request.UserId}
	} else if len(request.Nickname) > 0 {
		filter.Nicknames = []string{request.Nickname}
	} else {
		return resp, usecases.ErrInvalidRequest
	}

	to, err := c.usersStore.Find(ctx, filter)
	if err != nil {
		return resp, err
	}

	toUid := models.NoUser
	if len(to) > 0 {
		toUid = to[0].UserId
	}
	selfUid := uid

	directId, err := c.channelsStore.GetDirect(ctx, selfUid, toUid)
	if err != nil && err != store.ErrChanNotFound {
		return resp, err
	}
	if directId == models.NoChannel {
		directId, err = c.channelsStore.CreateDirect(ctx, selfUid, toUid)
		if err != nil {
			return resp, err
		}
	}
	direct, err := c.channelsStore.Get(ctx, directId)
	if err != nil {
		return resp, err
	}
	direct.Title = to[0].Title
	resp.SetChannel(&direct)
	return
}
