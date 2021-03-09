package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) Leave(ctx context.Context, uid models.Uid, request ChannelsLeaveRequest) (resp ChannelsOneResult, err error) {
	if err := c.membersStore.Leave(ctx, request.ChannelId, request.UserId); err != nil {
		resp.SetFailed(err)
		return resp, err
	}
	return c.getChan(ctx, request.ChannelId)
}
