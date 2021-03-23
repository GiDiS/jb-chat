package memory

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"sync"
	"time"
)

type sessionsMemoryStore struct {
	sessions map[string]models.Session
	rwMx     sync.RWMutex
}

func NewSessionsMemoryStore() *sessionsMemoryStore {
	return &sessionsMemoryStore{
		sessions: make(map[string]models.Session),
	}
}

func (s *sessionsMemoryStore) GetUserSessions(ctx context.Context, uid models.Uid) ([]models.Session, error) {
	list := make([]models.Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		if sess.UserId == uid {
			list = append(list, sess)
		}
	}
	return list, nil
}

func (s *sessionsMemoryStore) GetSessions(ctx context.Context) ([]models.Session, error) {
	list := make([]models.Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		list = append(list, sess)
	}
	return list, nil
}

func (s *sessionsMemoryStore) GetSession(ctx context.Context, sid string) (models.Session, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	sess, ok := s.sessions[sid]
	if ok {
		return sess, nil
	} else {
		now := time.Now()
		return models.Session{
			SessionId: sid,
			Started:   now,
			Updated:   now,
			IsOnline:  false,
		}, nil
	}
}

func (s *sessionsMemoryStore) SetSession(ctx context.Context, sid string, sess models.Session) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	s.sessions[sid] = sess
	return nil
}

func (s *sessionsMemoryStore) ClearSession(ctx context.Context, sid string) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	delete(s.sessions, sid)
	return nil
}

func (s *sessionsMemoryStore) SetOnline(ctx context.Context, sid string, uid models.Uid, isOnline bool) error {
	sess, err := s.GetSession(ctx, sid)
	if err != nil {
		return err
	}
	if isOnline {
		sess.Updated = time.Now()
		sess.IsOnline = true
	} else {
		sess.IsOnline = false
	}
	if err := s.SetSession(ctx, sid, sess); err != nil {
		return err
	}
	return nil
}

func (s *sessionsMemoryStore) IsOnline(ctx context.Context, uid models.Uid) (bool, error) {
	if uid == models.NoUser {
		return false, nil
	}
	userSessions, err := s.GetUserSessions(ctx, uid)
	if err != nil {
		return false, err
	}
	for _, sess := range userSessions {
		if sess.IsOnline {
			return true, nil
		}
	}
	return false, nil
}

func (s *sessionsMemoryStore) GetOnline(ctx context.Context, uid models.Uid) ([]string, error) {
	if uid == models.NoUser {
		return nil, nil
	}
	userSessions, err := s.GetUserSessions(ctx, uid)
	if err != nil {
		return nil, err
	}
	onlineSessions := make([]string, 0)
	for _, sess := range userSessions {
		if sess.IsOnline {
			onlineSessions = append(onlineSessions, sess.SessionId)
		}
	}
	return onlineSessions, nil
}

func (s *sessionsMemoryStore) GetAllOnline(ctx context.Context) (map[models.Uid][]string, error) {
	onlineSessions := make(map[models.Uid][]string)
	for _, sess := range s.sessions {
		if !sess.IsOnline {
			continue
		}
		if _, ok := onlineSessions[sess.UserId]; !ok {
			onlineSessions[sess.UserId] = make([]string, 1)
		}
		onlineSessions[sess.UserId] = append(onlineSessions[sess.UserId], sess.SessionId)
	}
	return onlineSessions, nil
}
