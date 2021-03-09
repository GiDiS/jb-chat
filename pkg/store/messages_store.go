package store

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
)

type MessagesStore interface {
	Create(context.Context, models.Message) (models.MessageId, error)
	Delete(context.Context, models.MessageId) error
	MarkAsThread(context.Context, models.MessageId, bool) error
	Get(context.Context, models.MessageId) (models.Message, error)
	Find(context.Context, MessagesSearchCriteria) ([]models.Message, error)
	Estimate(context.Context, MessagesSearchCriteria) (uint64, error)
}

type MessagesSearchCriteria struct {
	Ids       []models.MessageId `json:"ids"`
	ChannelId models.ChannelId   `json:"cid,omitempty"`
	ParentId  models.MessageId   `json:"pid,omitempty"`
	Search    string             `json:"search,omitempty"`
	Limits    models.Limits      `json:"limits,omitempty"`
}
