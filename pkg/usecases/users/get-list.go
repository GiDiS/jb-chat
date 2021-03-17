package users

import (
	"context"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (u *usersImpl) GetList(ctx context.Context, req UsersListRequest) (resp UsersListResponse, err error) {
	users, err := u.usersStore.Find(ctx, req.UserSearchCriteria)
	if err != nil {
		resp.SetFailed(err)
		return resp, fmt.Errorf("get user list failed: %w", err)
	}
	infos := make([]models.UserInfo, 0, len(users))
	for _, u := range users {
		infos = append(infos, models.UserInfo{
			User: u, Status: models.UserStatusOnline,
		})
	}
	resp.SetUsers(infos)
	return
}
