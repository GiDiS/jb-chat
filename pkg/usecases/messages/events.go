package messages

import "github.com/GiDiS/jb-chat/pkg/events"

const (
	MessageCreate  events.Type = "messages.create"
	MessageCreated events.Type = "messages.created"
	MessageUpdate  events.Type = "messages.update"
	MessageUpdated events.Type = "messages.updated"
	MessageDelete  events.Type = "messages.delete"
	MessageDeleted events.Type = "messages.deleted"
	MessageLike    events.Type = "messages.like"
	MessageUnlike  events.Type = "messages.unlike"
	MessageGetList events.Type = "messages.get-list"
	MessageList    events.Type = "messages.list"
	MessageGetInfo events.Type = "messages.get-info"
	MessageInfo    events.Type = "messages.info"
)
