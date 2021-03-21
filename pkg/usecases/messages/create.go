package messages

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"time"
)

func (c *messagesImpl) Create(ctx context.Context, request MessageCreateRequest) (resp MessageOneResult, err error) {

	// check channel exists
	if _, err := c.channelsStore.Get(ctx, request.ChannelId); err != nil {
		resp.SetFailed(err)
		return resp, err
	}

	msg := models.Message{
		ChannelId: request.ChannelId, UserId: request.UserId, Body: request.Body,
		Created: time.Now(),
	}

	if request.ParentId != nil {
		msg.ParentId = *request.ParentId
	}

	mid, err := c.messagesStore.Create(ctx, msg)
	if err != nil {
		resp.SetFailed(err)
		return resp, err
	}

	if err = c.channelsStore.SetLastMessage(ctx, msg.ChannelId, mid, msg.Created); err != nil {
		resp.SetFailed(err)
		return resp, err
	}

	return c.Get(ctx, mid)
}
