package container

import (
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/handlers_ws"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/usecases/auth"
	"github.com/GiDiS/jb-chat/pkg/usecases/channels"
	"github.com/GiDiS/jb-chat/pkg/usecases/messages"
	"github.com/GiDiS/jb-chat/pkg/usecases/system"
	"github.com/GiDiS/jb-chat/pkg/usecases/users"
)

const (
	Ping      events.Type = "ping"
	Pong      events.Type = "pong"
	Exit      events.Type = "exit"
	Broadcast events.Type = "broadcast"
)

var log logger.Logger

func GetEvents() map[events.Type]interface{} {
	return map[events.Type]interface{}{
		Ping:                       nil,
		Pong:                       nil,
		Broadcast:                  SysBroadcastPayload{},
		handlers_ws.WsConnected:    handlers_ws.SysClientResponse{},
		handlers_ws.WsDisconnected: handlers_ws.SysClientResponse{},

		system.SysGetConfig: nil,
		system.SysConfig:    system.ConfigResponse{},

		auth.AuthMe:         nil,
		auth.AuthMeInfo:     auth.AuthMeResponse{},
		auth.AuthRegister:   auth.AuthRegisterRequest{},
		auth.AuthRegistered: auth.AuthSignInResponse{},
		auth.AuthRequired:   nil,
		auth.AuthSignIn:     auth.AuthSignInRequest{},
		auth.AuthSignedIn:   auth.AuthSignInResponse{},
		auth.AuthSignOut:    nil,
		auth.AuthSignedOut:  auth.AuthSignOutResponse{},

		channels.ChannelsCreate:      channels.ChannelsCreateRequest{},
		channels.ChannelsCreated:     channels.ChannelsOneResult{},
		channels.ChannelsDelete:      channels.ChannelsDeleteRequest{},
		channels.ChannelsDeleted:     channels.ChannelsOneResult{},
		channels.ChannelsUpdate:      channels.ChannelsUpdateRequest{},
		channels.ChannelsUpdated:     channels.ChannelsOneResult{},
		channels.ChannelsJoin:        channels.ChannelsJoinRequest{},
		channels.ChannelsJoined:      channels.ChannelsOneResult{},
		channels.ChannelsLeave:       channels.ChannelsLeaveRequest{},
		channels.ChannelsLeft:        channels.ChannelsOneResult{},
		channels.ChannelsKick:        channels.ChannelsKickRequest{},
		channels.ChannelsKicked:      channels.ChannelsOneResult{},
		channels.ChannelsGetList:     channels.ChannelsGetListRequest{},
		channels.ChannelsList:        channels.ChannelsListResponse{},
		channels.ChannelsGetInfo:     channels.ChannelsGetInfoRequest{},
		channels.ChannelsInfo:        channels.ChannelsOneResult{},
		channels.ChannelsGetMembers:  channels.ChannelsMembersRequest{},
		channels.ChannelsMembers:     channels.ChannelsMembersResponse{},
		channels.ChannelsGetDirect:   channels.ChannelsGetDirectRequest{},
		channels.ChannelsDirectInfo:  channels.ChannelsOneResult{},
		channels.ChannelsGetLastSeen: channels.ChannelsGetLastSeenRequest{},
		channels.ChannelsLastSeen:    channels.ChannelsLastSeenResponse{},
		channels.ChannelsSetLastSeen: channels.ChannelsSetLastSeenRequest{},

		messages.MessageCreate:  messages.MessageCreateRequest{},
		messages.MessageCreated: messages.MessageOneResult{},
		messages.MessageUpdate:  messages.MessageUpdateRequest{},
		messages.MessageUpdated: messages.MessageOneResult{},
		messages.MessageDelete:  messages.MessageRef{},
		messages.MessageDeleted: messages.MessageOneResult{},
		messages.MessageLike:    messages.MessageRef{},
		messages.MessageUnlike:  messages.MessageRef{},
		messages.MessageGetList: messages.MessageListRequest{},
		messages.MessageList:    messages.MessageListResponse{},
		messages.MessageGetInfo: messages.MessageRef{},
		messages.MessageInfo:    messages.MessageOneResult{},

		users.UsersGetList: users.UsersListRequest{},
		users.UsersList:    users.UsersListResponse{},
		users.UsersGetInfo: users.UsersInfoRequest{},
		users.UsersInfo:    users.UsersInfoResponse{},
	}
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

type SysBroadcastPayload struct {
	Message string `json:"msg"`
}
