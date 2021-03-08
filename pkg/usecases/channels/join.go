package channels

import (
	"context"
	"jb_chat/pkg/models"
)

func (c *channelsImpl) Join(ctx context.Context, uid models.Uid, request ChannelsJoinRequest) (resp ChannelsOneResult, err error) {
	if err := c.membersStore.Join(ctx, request.ChannelId, request.UserId); err != nil {
		resp.SetFailed(err)
		return resp, err
	}
	return c.getChan(ctx, request.ChannelId)
}
