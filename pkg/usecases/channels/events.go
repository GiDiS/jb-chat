package channels

import "github.com/GiDiS/jb-chat/pkg/events"

const (
	ChannelsCreate     events.Type = "channels.create" // +
	ChannelsCreated    events.Type = "channels.created"
	ChannelsUpdate     events.Type = "channels.update"
	ChannelsUpdated    events.Type = "channels.updated"
	ChannelsDelete     events.Type = "channels.delete"
	ChannelsDeleted    events.Type = "channels.deleted"
	ChannelsJoin       events.Type = "channels.join" // +
	ChannelsJoined     events.Type = "channels.joined"
	ChannelsLeave      events.Type = "channels.leave" // +
	ChannelsLeft       events.Type = "channels.left"
	ChannelsKick       events.Type = "channels.kick"
	ChannelsKicked     events.Type = "channels.kicked"
	ChannelsGetList    events.Type = "channels.get-list"
	ChannelsList       events.Type = "channels.list"
	ChannelsGetInfo    events.Type = "channels.get-info" // +
	ChannelsInfo       events.Type = "channels.info"
	ChannelsGetMembers events.Type = "channels.get-members" // +
	ChannelsMembers    events.Type = "channels.members"
	ChannelsGetDirect  events.Type = "channels.get-direct" // +
	ChannelsDirectInfo events.Type = "channels.direct"
)
