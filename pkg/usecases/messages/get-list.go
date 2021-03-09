package messages

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
)

func (c *messagesImpl) GetList(ctx context.Context, request MessageListRequest) (resp MessageListResponse, err error) {
	filter := store.MessagesSearchCriteria{}
	filter.ChannelId = request.ChannelId
	if request.ParentId != nil {
		filter.ParentId = *request.ParentId
	}

	messages, _ := c.messagesStore.Find(ctx, filter)
	users := make([]models.UserInfo, 0, len(messages))
	usersMatched := make(map[models.Uid]bool)
	for _, m := range messages {
		if _, ok := usersMatched[m.UserId]; !ok {
			user, err := c.usersStore.GetByUid(ctx, m.UserId)
			if err != nil {
				continue
			}
			users = append(users, models.UserInfo{
				User: user, Status: models.UserStatusOnline,
			})
		}
	}
	resp.SetResult(messages, users)
	return
}
