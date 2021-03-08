package sessions

import (
	"context"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
)

type Sessions interface {
	Get(ctx context.Context, sid string) (models.Session, error)
	Reset(ctx context.Context, sid string) error
	Update(ctx context.Context, sid string, updater func(sess models.Session) (models.Session, error)) (models.Session, error)
	SetUid(ctx context.Context, sid string, uid models.Uid) error
	GetUid(ctx context.Context, sid string) (models.Uid, error)
}

type sessImpl struct {
	logger        logger.Logger
	sessionsStore store.SessionsStore
}

func NewSessions(logger logger.Logger, sessionsStore store.SessionsStore) *sessImpl {
	return &sessImpl{logger: logger, sessionsStore: sessionsStore}
}

func (s *sessImpl) Get(ctx context.Context, sid string) (sess models.Session, err error) {
	return s.sessionsStore.GetSession(sid)
}

func (s *sessImpl) Reset(ctx context.Context, sid string) error {
	return s.sessionsStore.ClearSession(sid)
}

func (s *sessImpl) Update(ctx context.Context, sid string, updater func(sess models.Session) (models.Session, error)) (models.Session, error) {
	sess, err := s.sessionsStore.GetSession(sid)
	if err != nil {
		return sess, err
	}
	sess, err = updater(sess)
	if err != nil {
		return sess, err
	}

	err = s.sessionsStore.SetSession(sid, sess)

	return sess, err
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
