package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) SetLastSeen(ctx context.Context, cid models.ChannelId, uid models.Uid, mid models.MessageId) error {
	return c.membersStore.SetLastSeen(ctx, cid, uid, mid)

}
