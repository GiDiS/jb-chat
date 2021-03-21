package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) GetLastSeen(ctx context.Context, cid models.ChannelId, uid models.Uid) (models.MessageId, error) {
	return c.membersStore.GetLastSeen(ctx, cid, uid)

}
