package postgres

import (
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db       *sqlx.DB
	channels *channelsPostgresStore
	messages *messagesPostgresStore
	sessions *sessionsPostgresStore
	users    *usersPostgresStore
}

func NewAppStore(db *sqlx.DB) (*PostgresStore, error) {
	return &PostgresStore{
		db:       db,
		channels: NewChannelsPostgresStore(db),
		messages: NewMessagesPostgresStore(db),
		sessions: NewSessionsPostgresStore(db),
		users:    NewUsersPostgresStore(db),
	}, nil
}

func (s PostgresStore) Channels() store.ChannelsStore {
	return s.channels
}

func (s PostgresStore) Members() store.ChannelMembersStore {
	return s.channels
}

func (s PostgresStore) Messages() store.MessagesStore {
	return s.messages
}

func (s PostgresStore) Sessions() store.SessionsStore {
	return s.sessions
}

func (s PostgresStore) OnlineUsers() store.UsersOnlineStore {
	return s.sessions
}

func (s PostgresStore) Users() store.UsersStore {
	return s.users
}
