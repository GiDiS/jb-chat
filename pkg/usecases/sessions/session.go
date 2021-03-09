package sessions

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

type Session interface {
	SetUid(ctx context.Context, sid string, uid models.Uid) error
	GetUid(ctx context.Context, sid string) (models.Uid, error)
}

func (s *sessImpl) SetUid(ctx context.Context, sid string, uid models.Uid) error {
	if sid == "" {
		return nil
	}
	_, err := s.Update(ctx, sid, func(sess models.Session) (models.Session, error) {
		sess.UserId = uid
		return sess, nil
	})

	return err
}

func (s *sessImpl) GetUid(ctx context.Context, sid string) (models.Uid, error) {
	if sid == "" {
		return models.NoUser, nil
	}
	sess, err := s.Get(ctx, sid)
	if err != nil {
		return models.NoUser, err
	}

	return sess.UserId, nil
}
