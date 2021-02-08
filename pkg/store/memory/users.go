package memory

import (
	"context"
	"jb-chat/pkg/models"
	"jb-chat/pkg/store"
	"sync"
)

type usersMemoryStore struct {
	users       map[models.Uid]models.User
	usersStatus map[models.Uid]models.UserStatus
	lastUid     models.Uid
	rwMx        sync.RWMutex
}

func NewUsersMemoryStore() *usersMemoryStore {
	return &usersMemoryStore{
		users:       make(map[models.Uid]models.User, 0),
		usersStatus: make(map[models.Uid]models.UserStatus, 0),
	}
}

func (s *usersMemoryStore) Register(ctx context.Context, user models.User) (models.Uid, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	existed, err := s.getByEmail(ctx, user.Email)
	if err == nil {
		return existed.Uid, store.ErrUserAlreadyRegistered
	} else if err != store.ErrUserNotFound {
		return models.NoUser, err
	}
	s.lastUid++
	user.Uid = s.lastUid
	s.users[user.Uid] = user
	return user.Uid, nil
}

func (s *usersMemoryStore) SetStatus(_ context.Context, uid models.Uid, status models.UserStatus) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	_, ok := s.users[uid]
	if !ok {
		return store.ErrUserNotFound
	}
	s.usersStatus[uid] = status
	return nil
}

func (s *usersMemoryStore) GetByEmail(ctx context.Context, email string) (models.User, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	return s.getByEmail(ctx, email)
}

func (s *usersMemoryStore) getByEmail(ctx context.Context, email string) (models.User, error) {
	return s.findOne(ctx, store.UserSearchCriteria{Emails: []string{email}})
}

func (s *usersMemoryStore) GetByUid(ctx context.Context, uid models.Uid) (models.User, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	return s.getByUid(ctx, uid)
}

func (s *usersMemoryStore) getByUid(ctx context.Context, uid models.Uid) (models.User, error) {
	return s.findOne(ctx, store.UserSearchCriteria{Uids: []models.Uid{uid}})
}

func (s *usersMemoryStore) Find(ctx context.Context, filter store.UserSearchCriteria) ([]models.User, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	return s.find(ctx, filter)
}

func (s *usersMemoryStore) FindActive(ctx context.Context, limits models.Limits) ([]models.User, error) {
	return s.Find(ctx, store.UserSearchCriteria{
		Statuses: []models.UserStatus{models.UserStatusOnline, models.UserStatusAway},
		Limits:   limits,
	})
}

func (s *usersMemoryStore) Estimate(_ context.Context, filter store.UserSearchCriteria) (uint64, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	var matched uint64 = 0
	for _, user := range s.users {
		if s.matchUser(&user, filter) {
			matched++
		}
	}
	return matched, nil
}

func (s *usersMemoryStore) find(_ context.Context, filter store.UserSearchCriteria) ([]models.User, error) {
	filtered := make([]models.User, 0)
	offset, limit, skipped := filter.Limits.Offset, filter.Limits.Limit, 0
	for _, user := range s.users {
		if s.matchUser(&user, filter) {
			if offset > 0 && skipped < offset {
				skipped++
				continue
			}
			filtered = append(filtered, user)

			if limit > 0 && len(filtered) >= limit {
				break
			}
		}
	}
	return filtered, nil
}

func (s *usersMemoryStore) findOne(ctx context.Context, filter store.UserSearchCriteria) (models.User, error) {
	matched, err := s.find(ctx, filter)
	if err != nil {
		return models.User{}, err
	}
	if len(matched) > 0 {
		return matched[0], nil
	}
	return models.User{}, store.ErrUserNotFound
}

func (s *usersMemoryStore) matchUser(user *models.User, filter store.UserSearchCriteria) bool {
	if user == nil {
		return false
	}

	if len(filter.Uids) > 0 {
		matched := false
		for _, fUid := range filter.Uids {
			if fUid != user.Uid {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(filter.Emails) > 0 {
		matched := false
		for _, fEmail := range filter.Emails {
			if fEmail != user.Email {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(filter.Statuses) > 0 {
		matched := false
		status := models.UserStatusUnknown
		if _, ok := s.usersStatus[user.Uid]; ok {
			status = s.usersStatus[user.Uid]
		}

		for _, fStatus := range filter.Statuses {
			if fStatus == status {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
