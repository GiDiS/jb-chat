package users

import (
	"context"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/store"
)

type Users interface {
	GetList(ctx context.Context, req UsersListRequest) (UsersListResponse, error)
	Get(ctx context.Context, req UsersInfoRequest) (UsersInfoResponse, error)
}

type usersImpl struct {
	logger     logger.Logger
	usersStore store.UsersStore
}

func NewUsers(logger logger.Logger, usersStore store.UsersStore) *usersImpl {
	return &usersImpl{logger: logger, usersStore: usersStore}
}
