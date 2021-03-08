package channels

import (
	"context"
	"jb_chat/pkg/models"
)

func (c *channelsImpl) Get(ctx context.Context, cid models.ChannelId) (ChannelsOneResult, error) {
	return c.getChan(ctx, cid)
}

func (c *channelsImpl) getChan(ctx context.Context, cid models.ChannelId) (ChannelsOneResult, error) {
	resp := ChannelsOneResult{}
	channel, err := c.channelsStore.Get(ctx, cid)
	if err != nil {
		resp.SetFailed(err)
		return resp, err
	}
	resp.SetChannel(&channel)
	return resp, nil
}
