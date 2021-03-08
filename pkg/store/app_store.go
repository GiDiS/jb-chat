package store

type AppStore interface {
	Channels() ChannelsStore
	Members() ChannelMembersStore
	Messages() MessagesStore
	Sessions() SessionsStore
	OnlineUsers() UsersOnlineStore
	Users() UsersStore
}
