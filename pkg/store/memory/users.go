package memory

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"sort"
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
		users:       make(map[models.Uid]models.User),
		usersStatus: make(map[models.Uid]models.UserStatus),
	}
}

func (s *usersMemoryStore) Register(ctx context.Context, user models.User) (models.Uid, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	return s.register(ctx, user)

}

func (s *usersMemoryStore) Save(ctx context.Context, user models.User) (models.Uid, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	if user.UserId > 0 {
		s.users[user.UserId] = user
		return user.UserId, nil
	} else {
		return s.register(ctx, user)
	}
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

func (s *usersMemoryStore) register(ctx context.Context, user models.User) (models.Uid, error) {
	existed, err := s.getByEmail(ctx, user.Email)
	if err == nil {
		return existed.UserId, store.ErrUserAlreadyRegistered
	} else if err != store.ErrUserNotFound {
		return models.NoUser, err
	}

	user.UserId = s.getNextUid()
	s.users[user.UserId] = user
	return user.UserId, nil
}

func (s *usersMemoryStore) getNextUid() models.Uid {
	lastUid := models.NoUser
	for _, u := range s.users {
		if lastUid < u.UserId {
			lastUid = u.UserId
		}
	}
	return lastUid + 1
}

func (s *usersMemoryStore) find(_ context.Context, filter store.UserSearchCriteria) ([]models.User, error) {
	filtered := make([]models.User, 0)
	offset, limit, skipped := filter.Limits.Offset, filter.Limits.Limit, 0
	for _, user := range s.users {
		if !s.matchUser(&user, filter) {
			continue
		}
		if offset > 0 && skipped < offset {
			skipped++
			continue
		}
		filtered = append(filtered, user)

		if limit > 0 && len(filtered) >= limit {
			break
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

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

	if filter.WithAvatars && user.AvatarUrl == "" {
		return false
	}

	if len(filter.Uids) > 0 {
		matched := false
		for _, fUid := range filter.Uids {
			if fUid == user.UserId {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(filter.Nicknames) > 0 {
		matched := false
		for _, fNickname := range filter.Nicknames {
			if fNickname == user.Nickname {
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
			if fEmail == user.Email {
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
		if _, ok := s.usersStatus[user.UserId]; ok {
			status = s.usersStatus[user.UserId]
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
