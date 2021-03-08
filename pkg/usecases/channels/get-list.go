package channels

import (
	"context"
)

func (c *channelsImpl) GetList(ctx context.Context, request ChannelsGetListRequest) (resp ChannelsListResponse, err error) {
	channels, err := c.channelsStore.Find(ctx, request.ChannelsSearchCriteria)
	if err != nil {
		resp.SetFailed(err)
	} else {
		resp.SetChannels(channels)
	}
	return
}
