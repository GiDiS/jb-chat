package store

import (
	"context"
	"jb_chat/pkg/models"
)

type SessionsStore interface {
	GetUserSessions(ctx context.Context, uid models.Uid) ([]models.Session, error)
	GetSessions(ctx context.Context) ([]models.Session, error)
	GetSession(ctx context.Context, sid string) (models.Session, error)
	SetSession(ctx context.Context, sid string, sess models.Session) error
	ClearSession(ctx context.Context, sid string) error
}

type UsersOnlineStore interface {
	SetOnline(ctx context.Context, sid string, uid models.Uid, isOnline bool) error
	IsOnline(ctx context.Context, uid models.Uid) (bool, error)
	GetOnline(ctx context.Context, uid models.Uid) ([]string, error)
	GetAllOnline(ctx context.Context) (map[models.Uid][]string, error)
}
