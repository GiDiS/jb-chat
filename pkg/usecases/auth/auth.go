package auth

import (
	"jb_chat/pkg/logger"
	"jb_chat/pkg/store"
)

type Auth interface {
	SignIn
	SignOut
}

type authImpl struct {
	logger     logger.Logger
	usersStore store.UsersStore
}

func NewAuth(logger logger.Logger, usersStore store.UsersStore) *authImpl {
	return &authImpl{logger: logger, usersStore: usersStore}
}
