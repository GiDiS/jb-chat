package memory

import (
	"jb_chat/pkg/models"
	"sync"
	"time"
)

type sessionsMemoryStore struct {
	sessions map[string]models.Session
	rwMx     sync.RWMutex
}

func NewSessionsMemoryStore() *sessionsMemoryStore {
	return &sessionsMemoryStore{
		sessions: make(map[string]models.Session, 0),
	}
}

func (s *sessionsMemoryStore) GetUserSessions(uid models.Uid) ([]models.Session, error) {
	list := make([]models.Session, len(s.sessions))
	for _, sess := range s.sessions {
		if sess.UserId == uid {
			list = append(list, sess)
		}
	}
	return list, nil
}

func (s *sessionsMemoryStore) GetSessions() ([]models.Session, error) {
	list := make([]models.Session, len(s.sessions))
	for _, sess := range s.sessions {
		list = append(list, sess)
	}
	return list, nil
}

func (s *sessionsMemoryStore) GetSession(sid string) (models.Session, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	sess, ok := s.sessions[sid]
	if ok {
		return sess, nil
	} else {
		return models.Session{SessionId: sid, Started: time.Now()}, nil
	}
}

func (s *sessionsMemoryStore) SetSession(sid string, sess models.Session) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	s.sessions[sid] = sess
	return nil
}

func (s *sessionsMemoryStore) ClearSession(sid string) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	delete(s.sessions, sid)
	return nil
}
