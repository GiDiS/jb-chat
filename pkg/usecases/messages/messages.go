package messages

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
)

type Messages interface {
	GetList(ctx context.Context, request MessageListRequest) (resp MessageListResponse, err error)
	Get(ctx context.Context, mid models.MessageId) (resp MessageOneResult, err error)
	Create(ctx context.Context, request MessageCreateRequest) (resp MessageOneResult, err error)
}

type messagesImpl struct {
	logger        logger.Logger
	messagesStore store.MessagesStore
	usersStore    store.UsersStore
}

func NewMessages(
	logger logger.Logger,
	messagesStore store.MessagesStore,
	usersStore store.UsersStore,
) *messagesImpl {
	return &messagesImpl{
		logger:        logger,
		messagesStore: messagesStore,
		usersStore:    usersStore,
	}
}
