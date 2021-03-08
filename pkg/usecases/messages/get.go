package messages

import (
	"context"
	"jb_chat/pkg/models"
)

func (c *messagesImpl) Get(ctx context.Context, mid models.MessageId) (resp MessageOneResult, err error) {

	if msg, err := c.messagesStore.Get(ctx, mid); err != nil {
		resp.SetFailed(err)
		return resp, err
	} else {
		resp.SetMsg(&msg)
	}
	return
}
