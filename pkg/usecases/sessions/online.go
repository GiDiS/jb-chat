package sessions

import (
	"context"
	"jb_chat/pkg/models"
)

type Online interface {
	SetOnline(ctx context.Context, sid string, uid models.Uid, isOnline bool) error
	IsOnline(ctx context.Context, uid models.Uid) (bool, error)
	GetOnline(ctx context.Context, uid models.Uid) ([]string, error)
	GetAllOnline(ctx context.Context) (map[models.Uid][]string, error)
}

func (s *sessImpl) SetOnline(ctx context.Context, sid string, uid models.Uid, isOnline bool) error {
	return s.usersOnlineStore.SetOnline(ctx, sid, uid, isOnline)
}

func (s *sessImpl) IsOnline(ctx context.Context, uid models.Uid) (bool, error) {
	return s.usersOnlineStore.IsOnline(ctx, uid)
}

func (s *sessImpl) GetOnline(ctx context.Context, uid models.Uid) ([]string, error) {
	return s.usersOnlineStore.GetOnline(ctx, uid)
}

func (s *sessImpl) GetAllOnline(ctx context.Context) (map[models.Uid][]string, error) {
	return s.usersOnlineStore.GetAllOnline(ctx)
}
