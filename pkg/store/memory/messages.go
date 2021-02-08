package memory

import (
	"context"
	"jb-chat/pkg/models"
	"jb-chat/pkg/store"
	"strings"
	"sync"
)

type messagesMemoryStore struct {
	messages        map[models.MessageId]models.Message
	channelMessages map[models.ChannelId]map[models.MessageId]bool
	lastMsg         models.MessageId

	rwMx sync.RWMutex
}

func NewMessagesMemoryStore() *messagesMemoryStore {
	return &messagesMemoryStore{}
}

func (s *messagesMemoryStore) Create(_ context.Context, message models.Message) (models.MessageId, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	if message.ChannelId <= 0 {
		return models.NoMessage, store.ErrEmptyChanId
	}

	if s.messages == nil {
		s.messages = make(map[models.MessageId]models.Message, 0)
	}
	s.lastMsg++
	message.MsgId = s.lastMsg
	s.messages[s.lastMsg] = message
	return message.MsgId, nil
}

func (s *messagesMemoryStore) Delete(ctx context.Context, id models.MessageId) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	msg, ok := s.messages[id]
	if !ok {
		return store.ErrMessageNotFound
	}

	// delete nested messages
	if msg.IsThread {
		nested, err := s.Find(ctx, store.MessagesSearchCriteria{ParentId: id})
		if err != nil {
			return err
		}
		for _, nMsg := range nested {
			delete(s.messages, nMsg.MsgId)
		}
	}

	delete(s.messages, id)

	return nil
}

func (s *messagesMemoryStore) MarkAsThread(_ context.Context, id models.MessageId, isThread bool) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	msg, ok := s.messages[id]
	if !ok {
		return store.ErrMessageNotFound
	}

	msg.IsThread = isThread
	s.messages[id] = msg

	return nil
}

func (s *messagesMemoryStore) Get(_ context.Context, id models.MessageId) (models.Message, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	if msg, ok := s.messages[id]; ok {
		return msg, nil
	} else {
		return models.Message{}, store.ErrMsgNotFound
	}
}

func (s *messagesMemoryStore) Find(_ context.Context, filter store.MessagesSearchCriteria) ([]models.Message, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	filtered := make([]models.Message, 0)
	offset, limit, skipped := filter.Limits.Offset, filter.Limits.Limit, 0
	for _, msg := range s.messages {
		if s.matchMessage(&msg, filter) {
			if offset > 0 && skipped < offset {
				skipped++
				continue
			}
			filtered = append(filtered, msg)

			if limit > 0 && len(filtered) >= limit {
				break
			}
		}
	}
	return filtered, nil
}

func (s *messagesMemoryStore) Estimate(_ context.Context, filter store.MessagesSearchCriteria) (uint64, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	var matched uint64 = 0
	for _, msg := range s.messages {
		if s.matchMessage(&msg, filter) {
			matched++
		}
	}
	return matched, nil
}

func (s *messagesMemoryStore) matchMessage(msg *models.Message, filter store.MessagesSearchCriteria) bool {
	if msg == nil {
		return false
	}

	if len(filter.Ids) > 0 {
		matched := false
		for _, id := range filter.Ids {
			if id == msg.MsgId {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if filter.ChannelId > 0 && msg.ChannelId != filter.ChannelId {
		return false
	}

	if filter.ParentId > 0 && msg.ParentId != filter.ParentId {
		return false
	}

	if filter.Search != "" && strings.Contains(msg.Body, filter.Search) {
		return false
	}

	return true
}
