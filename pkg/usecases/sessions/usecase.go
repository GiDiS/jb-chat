package sessions

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
)

type Sessions interface {
	Online
	Session
	Get(ctx context.Context, sid string) (models.Session, error)
	Reset(ctx context.Context, sid string) error
	Update(ctx context.Context, sid string, updater func(sess models.Session) (models.Session, error)) (models.Session, error)
}

type sessImpl struct {
	logger           logger.Logger
	sessionsStore    store.SessionsStore
	usersOnlineStore store.UsersOnlineStore
	sessions         map[string]models.Uid
}

func NewSessions(
	logger logger.Logger,
	sessionsStore store.SessionsStore,
	usersOnlineStore store.UsersOnlineStore,
) *sessImpl {
	return &sessImpl{
		logger:           logger,
		sessionsStore:    sessionsStore,
		usersOnlineStore: usersOnlineStore,
		sessions:         make(map[string]models.Uid, 0),
	}
}

func (s *sessImpl) Get(ctx context.Context, sid string) (sess models.Session, err error) {
	return s.sessionsStore.GetSession(ctx, sid)
}

func (s *sessImpl) Reset(ctx context.Context, sid string) error {
	return s.sessionsStore.ClearSession(ctx, sid)
}

func (s *sessImpl) Update(ctx context.Context, sid string, updater func(sess models.Session) (models.Session, error)) (models.Session, error) {
	sess, err := s.sessionsStore.GetSession(ctx, sid)
	if err != nil {
		return sess, err
	}
	sess, err = updater(sess)
	if err != nil {
		return sess, err
	}

	err = s.sessionsStore.SetSession(ctx, sid, sess)

	return sess, err
}
