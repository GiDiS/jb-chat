package memory

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"sort"
	"strings"
	"sync"
)

type channelsMemoryStore struct {
	channels      map[models.ChannelId]models.Channel
	messages      map[models.ChannelId][]models.Message
	members       map[models.ChannelId]map[models.Uid]bool
	usersChannels map[models.Uid]map[models.ChannelId]bool
	lastCid       models.ChannelId
	rwMx          sync.RWMutex
}

func NewChannelsMemoryStore() *channelsMemoryStore {
	return &channelsMemoryStore{
		usersChannels: make(map[models.Uid]map[models.ChannelId]bool),
		channels:      make(map[models.ChannelId]models.Channel),
		members:       make(map[models.ChannelId]map[models.Uid]bool),
	}
}

func (s *channelsMemoryStore) CreateDirect(ctx context.Context, uidA models.Uid, uidB models.Uid) (models.ChannelId, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	s.lastCid++
	cid := s.lastCid
	s.channels[cid] = models.Channel{
		Cid:          cid,
		Title:        models.DirectTitle(uidA, uidB),
		LastMsgId:    models.NoMessage,
		MembersCount: 1,
		Type:         models.ChannelTypeDirect,
	}

	if err := s.join(ctx, cid, uidA); err != nil {
		return 0, err
	}
	if err := s.join(ctx, cid, uidB); err != nil {
		return 0, err
	}

	return cid, nil
}

func (s *channelsMemoryStore) GetDirect(ctx context.Context, uidA models.Uid, uidB models.Uid) (models.ChannelId, error) {
	channels, err := s.Find(ctx, store.ChannelsSearchCriteria{
		Type:  models.ChannelTypeDirect,
		Title: models.DirectTitle(uidA, uidB),
	})
	if err != nil {
		return models.NoChannel, err
	}
	if len(channels) > 0 {
		return channels[0].Cid, nil
	}
	return models.NoChannel, nil
}

func (s *channelsMemoryStore) CreatePublic(ctx context.Context, authorUid models.Uid, title string) (models.ChannelId, error) {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	s.lastCid++
	cid := s.lastCid
	s.channels[cid] = models.Channel{
		Cid:          cid,
		Title:        title,
		LastMsgId:    models.NoMessage,
		MembersCount: 1,
		Type:         models.ChannelTypePublic,
	}

	if err := s.join(ctx, cid, authorUid); err != nil {
		return 0, err
	}

	return cid, nil
}

func (s *channelsMemoryStore) Delete(_ context.Context, cid models.ChannelId) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	delete(s.channels, cid)
	delete(s.members, cid)
	for uid, channels := range s.usersChannels {
		if _, ok := channels[cid]; ok {
			delete(channels, cid)
		}
		s.usersChannels[uid] = channels
	}

	return nil
}

func (s *channelsMemoryStore) Get(_ context.Context, cid models.ChannelId) (models.Channel, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	if ch, ok := s.channels[cid]; ok {
		return ch, nil
	}
	return models.Channel{}, store.ErrChanNotFound
}

func (s *channelsMemoryStore) Find(_ context.Context, filter store.ChannelsSearchCriteria) ([]models.Channel, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	filtered := make([]models.Channel, 0)
	offset, limit, skipped := filter.Limits.Offset, filter.Limits.Limit, 0
	for _, ch := range s.channels {
		if s.channelMatch(&ch, filter) {
			if offset > 0 && skipped < offset {
				skipped++
				continue
			}
			filtered = append(filtered, ch)

			if limit > 0 && len(filtered) >= limit {
				break
			}
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return strings.Compare(filtered[i].Title, filtered[j].Title) < 0
	})

	return filtered, nil
}

func (s *channelsMemoryStore) Estimate(_ context.Context, filter store.ChannelsSearchCriteria) (uint64, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	var total uint64 = 0
	for _, ch := range s.channels {
		if s.channelMatch(&ch, filter) {
			total++
		}
	}
	return total, nil
}

func (s *channelsMemoryStore) channelMatch(ch *models.Channel, filter store.ChannelsSearchCriteria) bool {
	if ch == nil {
		return false
	}

	if len(filter.ChannelIds) > 0 {
		match := false
		for _, cid := range filter.ChannelIds {
			if cid == ch.Cid {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if int(filter.HasMember) > 0 {
		match := false

		members, ok := s.members[ch.Cid]
		if !ok {
			return false
		}
		for uid, active := range members {
			if uid == filter.HasMember && active {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if filter.Title != "" && !strings.Contains(ch.Title, filter.Title) {
		return false
	}
	if filter.Type != models.ChannelTypeUnknown && filter.Type != ch.Type {
		return false
	}

	return true
}

func (s *channelsMemoryStore) Members(_ context.Context, cid models.ChannelId) ([]models.Uid, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	if _, ok := s.members[cid]; ok {
		members := make([]models.Uid, 0, len(s.members[cid]))
		for uid, active := range s.members[cid] {
			if active {
				members = append(members, uid)
			}
		}
		return members, nil
	}
	return nil, store.ErrChanNotFound
}

func (s *channelsMemoryStore) MemberOf(_ context.Context, uid models.Uid) ([]models.ChannelId, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()

	channels := make([]models.ChannelId, 0)
	for cid, active := range s.usersChannels[uid] {
		if active {
			channels = append(channels, cid)
		}
	}

	return channels, nil
}

func (s *channelsMemoryStore) IsMember(ctx context.Context, cid models.ChannelId, uid models.Uid) (bool, error) {
	s.rwMx.RLock()
	defer s.rwMx.RUnlock()
	return s.isMember(ctx, cid, uid)
}

func (s *channelsMemoryStore) isMember(_ context.Context, cid models.ChannelId, uid models.Uid) (bool, error) {
	members, ok := s.members[cid]
	if !ok {
		return false, nil
	}
	flag, ok := members[uid]
	if ok && flag {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *channelsMemoryStore) Join(ctx context.Context, cid models.ChannelId, uid models.Uid) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()
	return s.join(ctx, cid, uid)
}

func (s *channelsMemoryStore) join(ctx context.Context, cid models.ChannelId, uid models.Uid) error {
	if isMember, err := s.isMember(ctx, cid, uid); err != nil {
		return err
	} else if isMember {
		return nil
	}

	if _, ok := s.members[cid]; !ok {
		s.members[cid] = make(map[models.Uid]bool)
	}
	s.members[cid][uid] = true
	if _, ok := s.usersChannels[uid]; !ok {
		s.usersChannels[uid] = make(map[models.ChannelId]bool)
	}
	s.usersChannels[uid][cid] = true

	return nil
}

func (s *channelsMemoryStore) Leave(ctx context.Context, cid models.ChannelId, uid models.Uid) error {
	s.rwMx.Lock()
	defer s.rwMx.Unlock()

	if isMember, err := s.isMember(ctx, cid, uid); err != nil {
		return err
	} else if !isMember {
		return nil
	}

	delete(s.members[cid], uid)
	delete(s.usersChannels[uid], cid)

	return nil
}
