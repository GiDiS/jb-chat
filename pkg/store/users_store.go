package store

import (
	"context"
	"jb-chat/pkg/models"
)

type UsersStore interface {
	Register(context.Context, models.User) (models.Uid, error)
	SetStatus(context.Context, models.Uid, models.UserStatus) error
	GetByEmail(context.Context, string) (models.User, error)
	GetByUid(context.Context, models.Uid) (models.User, error)
	Find(context.Context, UserSearchCriteria) ([]models.User, error)
	Estimate(context.Context, UserSearchCriteria) (uint64, error)
}

type UserSearchCriteria struct {
	Uids     []models.Uid        `json:"uids,omitempty"`
	Emails   []string            `json:"emails,omitempty"`
	Statuses []models.UserStatus `json:"statuses,omitempty"`
	Limits   models.Limits       `json:"limits,omitempty"`
}
