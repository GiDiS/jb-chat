package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) Create(ctx context.Context, uid models.Uid, request ChannelsCreateRequest) (resp ChannelsOneResult, err error) {
	cid, err := c.channelsStore.CreatePublic(ctx, uid, request.Title)
	if err != nil {
		return resp, err
	}
	return c.getChan(ctx, cid)
}
