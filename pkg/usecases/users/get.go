package users

import (
	"context"
	"jb_chat/pkg/models"
)

func (u *usersImpl) Get(ctx context.Context, req UsersInfoRequest) (resp UsersInfoResponse, err error) {
	user, err := u.usersStore.GetByUid(ctx, req.UserId)
	if err != nil {
		resp.SetFailed(err)
		return resp, err
	} else {
		resp.SetUser(models.UserInfo{User: user})
	}
	return
}
