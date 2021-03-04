package store

type AppStore interface {
	Channels() ChannelsStore
	Members() ChannelMembersStore
	Messages() MessagesStore
	Sessions() SessionsStore
	Users() UsersStore
}
