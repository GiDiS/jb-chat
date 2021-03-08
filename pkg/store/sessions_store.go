package store

import "jb_chat/pkg/models"

type SessionsStore interface {
	GetUserSessions(uid models.Uid) ([]models.Session, error)
	GetSessions() ([]models.Session, error)
	GetSession(sid string) (models.Session, error)
	SetSession(sid string, sess models.Session) error
	ClearSession(sid string) error
}
