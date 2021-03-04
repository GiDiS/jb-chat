package event

import (
	"jb_chat/pkg/events"
	"jb_chat/pkg/logger"
)

const (
	Ping      events.Type = "ping"
	Pong      events.Type = "pong"
	Exit      events.Type = "exit"
	Broadcast events.Type = "broadcast"

	AuthMe         events.Type = "auth.me"
	AuthMeInfo     events.Type = "auth.me-info"
	AuthRegister   events.Type = "auth.register"
	AuthRegistered events.Type = "auth.registered"
	AuthRequired   events.Type = "auth.required"
	AuthSignIn     events.Type = "auth.sign-in"
	AuthSignedIn   events.Type = "auth.signed-in"
	AuthSignOut    events.Type = "auth.sign-out"
	AuthSignedOut  events.Type = "auth.signed-out"

	ChannelsCreate     events.Type = "channels.create"
	ChannelsCreated    events.Type = "channels.created"
	ChannelsUpdate     events.Type = "channels.update"
	ChannelsUpdated    events.Type = "channels.updated"
	ChannelsDelete     events.Type = "channels.delete"
	ChannelsDeleted    events.Type = "channels.deleted"
	ChannelsJoin       events.Type = "channels.join"
	ChannelsJoined     events.Type = "channels.joined"
	ChannelsLeave      events.Type = "channels.leave"
	ChannelsLeft       events.Type = "channels.left"
	ChannelsKick       events.Type = "channels.kick"
	ChannelsKicked     events.Type = "channels.kicked"
	ChannelsGetList    events.Type = "channels.get-list"
	ChannelsList       events.Type = "channels.list"
	ChannelsGetInfo    events.Type = "channels.get-info"
	ChannelsInfo       events.Type = "channels.info"
	ChannelsGetMembers events.Type = "channels.get-members"
	ChannelsMembers    events.Type = "channels.members"
	ChannelsGetDirect  events.Type = "channels.get-direct"
	ChannelsDirectInfo events.Type = "channels.direct"

	UsersGetList events.Type = "users.get-list"
	UsersList    events.Type = "users.list"
	UsersGetInfo events.Type = "users.get-info"
	UsersInfo    events.Type = "users.info"

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

var eventsMap = map[events.Type]interface{}{
	Ping:           nil,
	Pong:           nil,
	Broadcast:      SysBroadcastPayload{},
	AuthMe:         nil,
	AuthMeInfo:     AuthMeResponse{},
	AuthRegister:   AuthRegisterRequest{},
	AuthRegistered: AuthSignInResponse{},
	AuthRequired:   nil,
	AuthSignIn:     AuthSignInRequest{},
	AuthSignedIn:   AuthSignInResponse{},
	AuthSignOut:    nil,
	AuthSignedOut:  AuthSignOutResponse{},

	ChannelsCreate:     ChannelsCreateRequest{},
	ChannelsCreated:    ChannelsOneResult{},
	ChannelsDelete:     ChannelsDeleteRequest{},
	ChannelsDeleted:    ChannelsOneResult{},
	ChannelsUpdate:     ChannelsUpdateRequest{},
	ChannelsUpdated:    ChannelsOneResult{},
	ChannelsJoin:       ChannelsJoinRequest{},
	ChannelsJoined:     ChannelsOneResult{},
	ChannelsLeave:      ChannelsLeaveRequest{},
	ChannelsLeft:       ChannelsOneResult{},
	ChannelsKick:       ChannelsKickRequest{},
	ChannelsKicked:     ChannelsOneResult{},
	ChannelsGetList:    ChannelsGetListRequest{},
	ChannelsList:       ChannelsListResponse{},
	ChannelsGetInfo:    ChannelsGetInfoRequest{},
	ChannelsInfo:       ChannelsOneResult{},
	ChannelsGetMembers: ChannelsMembersRequest{},
	ChannelsMembers:    ChannelsMembersResponse{},
	ChannelsGetDirect:  ChannelsGetDirectRequest{},
	ChannelsDirectInfo: ChannelsOneResult{},

	MessageCreate:  MessageCreateRequest{},
	MessageCreated: MessageOneResult{},
	MessageUpdate:  MessageUpdateRequest{},
	MessageUpdated: MessageOneResult{},
	MessageDelete:  MessageRef{},
	MessageDeleted: MessageOneResult{},
	MessageLike:    MessageRef{},
	MessageUnlike:  MessageRef{},
	MessageGetList: MessageListRequest{},
	MessageList:    MessageListResponse{},
	MessageGetInfo: MessageRef{},
	MessageInfo:    MessageOneResult{},

	UsersGetList: UsersListRequest{},
	UsersList:    UsersListResponse{},
	UsersGetInfo: UsersInfoRequest{},
	UsersInfo:    UsersInfoResponse{},
}

var log logger.Logger

type ResultStatus struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func GetEvents() map[events.Type]interface{} {
	return eventsMap
}

func init() {
	log = logger.DefaultLogger()
	resolver := events.DefaultResolver
	for eventType, proto := range GetEvents() {
		if err := resolver.Register(eventType, proto); err != nil {
			log.WithError(err).Fatalf("register event failed")
		}
	}
}

func (r *ResultStatus) IsFailed() bool {
	return !r.Ok
}

func (r *ResultStatus) SetSuccess(msg string) {
	r.Ok = true
	r.Message = msg
}

func (r *ResultStatus) SetFailed(err error) {
	r.Ok = false
	r.Message = err.Error()
}
