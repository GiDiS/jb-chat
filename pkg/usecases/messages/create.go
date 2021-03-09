package messages

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"time"
)

func (c *messagesImpl) Create(ctx context.Context, request MessageCreateRequest) (resp MessageOneResult, err error) {
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

	return c.Get(ctx, mid)
}
