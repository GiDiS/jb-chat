package users

import (
	"context"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
)

func (u *usersImpl) Get(ctx context.Context, req UsersInfoRequest) (resp UsersInfoResponse, err error) {
	user, err := u.usersStore.GetByUid(ctx, req.UserId)
	if err != nil {
		resp.SetFailed(err)
		return resp, fmt.Errorf("get user failed: %w", err)
	}
	resp.SetUser(models.UserInfo{User: user})
	return
}
