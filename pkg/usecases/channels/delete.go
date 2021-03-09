package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) Delete(ctx context.Context, uid models.Uid, request ChannelsDeleteRequest) (resp ChannelsOneResult, err error) {
	// @todo check access
	err = c.channelsStore.Delete(ctx, request.ChannelId)
	if err != nil {
		resp.SetFailed(err)
	} else {
		resp.ChannelId = request.ChannelId
		resp.SetSuccess("ok")
	}
	return
}
