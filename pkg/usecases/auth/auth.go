package auth

import (
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/store"
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
