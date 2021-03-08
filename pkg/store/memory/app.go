package memory

import "jb_chat/pkg/store"

type memoryStore struct {
	channels *channelsMemoryStore
	messages *messagesMemoryStore
	sessions *sessionsMemoryStore
	users    *usersMemoryStore
}

func NewAppStore() *memoryStore {
	return &memoryStore{
		channels: NewChannelsMemoryStore(),
		messages: NewMessagesMemoryStore(),
		sessions: NewSessionsMemoryStore(),
		users:    NewUsersMemoryStore(),
	}
}

func (s memoryStore) Channels() store.ChannelsStore {
	return s.channels
}

func (s memoryStore) Members() store.ChannelMembersStore {
	return s.channels
}

func (s memoryStore) Messages() store.MessagesStore {
	return s.messages
}

func (s memoryStore) Sessions() store.SessionsStore {
	return s.sessions
}

func (s memoryStore) OnlineUsers() store.UsersOnlineStore {
	return s.sessions
}

func (s memoryStore) Users() store.UsersStore {
	return s.users
}
