package channels

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (c *channelsImpl) GetMembers(ctx context.Context, uid models.Uid, request ChannelsMembersRequest) (resp ChannelsMembersResponse, err error) {
	members, err := c.membersStore.Members(ctx, request.ChannelId)
	if err != nil {
		resp.SetFailed(err)
	} else {
		resp.SetMembers(members)
	}
	return
}
