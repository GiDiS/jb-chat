package container

import (
	"context"
	"errors"
	"jb_chat/pkg/events"
	"jb_chat/pkg/handlers_ws"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/models"
	"jb_chat/pkg/usecases"
	authUc "jb_chat/pkg/usecases/auth"
	channelsUc "jb_chat/pkg/usecases/channels"
	messagesUc "jb_chat/pkg/usecases/messages"
	sessionsUc "jb_chat/pkg/usecases/sessions"
	usersUc "jb_chat/pkg/usecases/users"
	"sync"
)

type Dispatcher struct {
	dispatcher events.Dispatcher
	logger     logger.Logger
	authUc     authUc.Auth
	channelsUc channelsUc.Channels
	messagesUc messagesUc.Messages
	sessionsUc sessionsUc.Sessions
	usersUc    usersUc.Users
	mx         sync.Mutex
}

func NewDispatcher(c Container) *Dispatcher {
	d := Dispatcher{
		dispatcher: c.EventsDispatcher,
		logger:     c.Logger,
		authUc:     authUc.NewAuth(c.Logger, c.Store.Users()),
		channelsUc: channelsUc.NewChannels(c.Logger, c.Store.Channels(), c.Store.Members(), c.Store.Users()),
		messagesUc: messagesUc.NewMessages(c.Logger, c.Store.Messages(), c.Store.Users()),
		sessionsUc: sessionsUc.NewSessions(c.Logger, c.Store.Sessions(), c.Store.OnlineUsers()),
		usersUc:    usersUc.NewUsers(c.Logger, c.Store.Users()),
	}
	d.init()
	return &d
}

func (d *Dispatcher) init() {
	d.on(Ping, d.onPing)
	d.on(Pong, d.onPong)
	d.on(handlers_ws.WsConnected, d.onConnected)
	d.on(handlers_ws.WsDisconnected, d.onDisconnected)
	d.on(Broadcast, d.onBroadcast)
	d.on(authUc.AuthRegister, d.onAuthRegister)
	d.on(authUc.AuthSignIn, d.onAuthSignIn)
	d.on(authUc.AuthSignOut, d.onAuthSignOut)

	d.onRegistered(channelsUc.ChannelsGetList, d.onChannelsGetList)
	d.onRegistered(channelsUc.ChannelsGetInfo, d.onChannelsGet)
	d.onRegistered(channelsUc.ChannelsGetDirect, d.onChannelsGetDirect)
	d.onRegistered(channelsUc.ChannelsGetMembers, d.onChannelsGetMembers)
	d.onRegistered(channelsUc.ChannelsCreate, d.onChannelsCreate)
	d.onRegistered(channelsUc.ChannelsDelete, d.onChannelsDelete)
	d.onRegistered(channelsUc.ChannelsJoin, d.onChannelsJoin)
	d.onRegistered(channelsUc.ChannelsLeave, d.onChannelsLeave)

	d.onRegistered(usersUc.UsersGetList, d.onUsersGetList)
	d.onRegistered(usersUc.UsersGetInfo, d.onUsersGetInfo)

	d.onRegistered(messagesUc.MessageGetList, d.onMessagesGetList)
	d.onRegistered(messagesUc.MessageCreate, d.onMessageCreate)
}

func (d *Dispatcher) onPing(e events.Event) error {
	d.reply(Pong, e, nil)
	return nil
}

func (d *Dispatcher) onPong(e events.Event) error {
	d.reply(Ping, e, nil)

	return nil
}

func (d *Dispatcher) onBroadcast(e events.Event) error {
	d.broadcast(Broadcast, e, e.Payload)
	return nil
}

func (d *Dispatcher) onConnected(e events.Event) error {
	d.logger.Debugf("Connected: %v", e.Payload)
	d.broadcast(handlers_ws.WsConnected, e, e.Payload)

	return nil
}

func (d *Dispatcher) onDisconnected(e events.Event) error {
	d.logger.Debugf("Disconnected: %v", e.Payload)

	sid, uid := e.GetSid(), e.GetUid()
	if err := d.sessionsUc.SetOnline(e.Ctx, sid, uid, false); err != nil {
		return err
	}
	if err := d.sessionsUc.Reset(e.Ctx, sid); err != nil {
		return err
	}

	d.broadcast(handlers_ws.WsDisconnected, e, e.Payload)
	d.broadcastUserInfo(e.Ctx, e, uid)

	return nil
}

func (d *Dispatcher) onAuthRegister(e events.Event) error {
	payload, ok := e.Payload.(authUc.AuthRegisterRequest)
	if !ok {
		return errors.New("wrong req")
	}
	uid := e.GetUid()
	if uid == models.NoUser {
		return errors.New("auth required")
	}

	d.logger.Debug(payload)

	d.reply(authUc.AuthRegistered, e, payload)
	d.broadcastUserInfo(e.Ctx, e, uid)

	return nil
}

func (d *Dispatcher) onAuthSignIn(e events.Event) error {
	payload, ok := e.Payload.(authUc.AuthSignInRequest)
	if !ok || e.Type != authUc.AuthSignIn {
		return usecases.ErrInvalidRequest
	}
	resp, err := d.authUc.SignIn(e.Ctx, payload)
	if err != nil {
		return err
	} else if resp.Me == nil {
		d.reply(authUc.AuthRequired, e, resp)
		return nil
	}

	sid := e.GetSid()
	if err = d.sessionsUc.SetUid(e.Ctx, sid, resp.Me.UserId); err != nil {
		return err
	}
	if err = d.sessionsUc.SetOnline(e.Ctx, sid, resp.Me.UserId, true); err != nil {
		return err
	}

	d.reply(authUc.AuthSignedIn, e, resp)
	d.broadcastUserInfo(e.Ctx, e, resp.Me.UserId)

	return nil
}

func (d *Dispatcher) onAuthSignOut(e events.Event) error {
	sid, uid := e.GetSid(), e.GetUid()
	if err := d.authUc.SignOut(e.Ctx, ""); err != nil {
		return err
	}

	if err := d.sessionsUc.SetOnline(e.Ctx, sid, uid, false); err != nil {
		return err
	}
	if err := d.sessionsUc.SetUid(e.Ctx, sid, models.NoUser); err != nil {
		return err
	}

	d.reply(authUc.AuthSignedOut, e, nil)
	d.broadcastUserInfo(e.Ctx, e, uid)

	return nil
}

func (d *Dispatcher) onChannelsGetList(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsGetListRequest)
	if e.Type != channelsUc.ChannelsGetList || !ok {
		return usecases.ErrInvalidRequest
	}
	if resp, err := d.channelsUc.GetList(e.Ctx, request); err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsList, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsGet(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsGetInfoRequest)
	if e.Type != channelsUc.ChannelsGetInfo || !ok {
		return usecases.ErrInvalidRequest
	}
	if resp, err := d.channelsUc.Get(e.Ctx, request.ChannelId); err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsInfo, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsGetDirect(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsGetDirectRequest)
	if e.Type != channelsUc.ChannelsGetDirect || !ok {
		return usecases.ErrInvalidRequest
	}
	if resp, err := d.channelsUc.GetDirect(e.Ctx, e.GetUid(), request); err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsDirectInfo, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsGetMembers(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsMembersRequest)
	if e.Type != channelsUc.ChannelsGetMembers || !ok {
		return usecases.ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	if resp, err := d.channelsUc.GetMembers(e.Ctx, e.GetUid(), request); err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsMembers, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsCreate(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsCreateRequest)
	if e.Type != channelsUc.ChannelsCreate || !ok {
		return usecases.ErrInvalidRequest
	}

	if resp, err := d.channelsUc.Create(e.Ctx, e.GetUid(), request); err != nil {
		return err
	} else {
		d.broadcast(channelsUc.ChannelsCreated, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsDelete(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsDeleteRequest)
	if e.Type != channelsUc.ChannelsDelete || !ok {
		return usecases.ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	if resp, err := d.channelsUc.Delete(e.Ctx, e.GetUid(), request); err != nil {
		return err
	} else {
		d.broadcast(channelsUc.ChannelsDeleted, e, resp)
	}
	return nil
}

func (d *Dispatcher) onChannelsJoin(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsJoinRequest)
	if e.Type != channelsUc.ChannelsJoin || !ok {
		return usecases.ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	if resp, err := d.channelsUc.Join(e.Ctx, e.GetUid(), request); err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsJoined, e, resp)
	}
	d.replyChannelMembers(e.Ctx, e, request.ChannelId)
	return nil
}

func (d *Dispatcher) onChannelsLeave(e events.Event) error {
	request, ok := e.Payload.(channelsUc.ChannelsLeaveRequest)
	if e.Type != channelsUc.ChannelsLeave || !ok {
		return usecases.ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	resp, err := d.channelsUc.Leave(e.Ctx, e.GetUid(), request)
	if err != nil {
		return err
	} else {
		d.reply(channelsUc.ChannelsLeft, e, resp)
	}
	d.replyChannelMembers(e.Ctx, e, request.ChannelId)
	return nil
}

func (d *Dispatcher) onUsersGetList(e events.Event) error {
	request, ok := e.Payload.(usersUc.UsersListRequest)
	if e.Type != usersUc.UsersGetList || !ok {
		return usecases.ErrInvalidRequest
	}

	resp, err := d.usersUc.GetList(e.Ctx, request)
	if err != nil {
		return err
	}
	online, err := d.sessionsUc.GetAllOnline(e.Ctx)
	if err != nil {
		return err
	}
	users := resp.Users
	for idx, user := range users {
		if _, ok := online[user.UserId]; ok {
			user.Status = models.UserStatusOnline
		} else {
			user.Status = models.UserStatusOffline
		}
		users[idx] = user
	}
	resp.SetUsers(users)

	d.reply(usersUc.UsersList, e, resp)

	return nil
}

func (d *Dispatcher) onUsersGetInfo(e events.Event) error {
	request, ok := e.Payload.(usersUc.UsersInfoRequest)
	if e.Type != usersUc.UsersGetInfo || !ok {
		return usecases.ErrInvalidRequest
	}

	if resp, err := d.usersUc.Get(e.Ctx, request); err != nil {
		return err
	} else {
		d.reply(usersUc.UsersInfo, e, resp)
	}
	return nil
}

func (d *Dispatcher) onMessagesGetList(e events.Event) error {
	request, ok := e.Payload.(messagesUc.MessageListRequest)
	if e.Type != messagesUc.MessageGetList || !ok {
		return usecases.ErrInvalidRequest
	}

	resp, err := d.messagesUc.GetList(e.Ctx, request)
	if err != nil {
		return err
	}

	online, err := d.sessionsUc.GetAllOnline(e.Ctx)
	if err != nil {
		return err
	}
	users := resp.Users
	for idx, user := range users {
		if _, ok := online[user.UserId]; ok {
			user.Status = models.UserStatusOnline
		} else {
			user.Status = models.UserStatusOffline
		}
		users[idx] = user
	}
	resp.Users = users
	d.reply(messagesUc.MessageList, e, resp)

	return nil
}

func (d *Dispatcher) onMessageCreate(e events.Event) error {
	request, ok := e.Payload.(messagesUc.MessageCreateRequest)
	if e.Type != messagesUc.MessageCreate || !ok {
		return usecases.ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	if resp, err := d.messagesUc.Create(e.Ctx, request); err != nil {
		return err
	} else {
		d.broadcast(messagesUc.MessageCreated, e, resp)
	}
	return nil
}

func (d *Dispatcher) replyChannelMembers(ctx context.Context, prev events.Event, cid models.ChannelId) {
	if cid == models.NoChannel {
		return
	}
	resp, err := d.channelsUc.GetMembers(ctx, prev.GetUid(), channelsUc.ChannelsMembersRequest{ChannelId: cid})
	if err != nil {
		d.logger.Errorf("Get members failed: %v", err)
		return
	}
	d.reply(channelsUc.ChannelsMembers, prev, resp)
}

func (d *Dispatcher) broadcastUserInfo(ctx context.Context, prev events.Event, uid models.Uid) {
	if uid == models.NoUser {
		return
	}
	resp, err := d.usersUc.Get(ctx, usersUc.UsersInfoRequest{
		UserId: uid,
	})
	if err != nil {
		d.logger.Errorf("users get failed: %v", err)
		return
	}

	isOnline, err := d.sessionsUc.IsOnline(ctx, uid)
	if err != nil {
		d.logger.Errorf("user get status failed: %v", err)
		resp.User.Status = models.UserStatusUnknown
	} else if isOnline {
		resp.User.Status = models.UserStatusOnline
	} else {
		resp.User.Status = models.UserStatusOffline
	}

	d.broadcast(usersUc.UsersInfo, prev, resp)
}

func (d *Dispatcher) on(t events.Type, h events.EventHandler) {
	d.dispatcher.On(t, d.auth(h, false))
}

func (d *Dispatcher) onRegistered(t events.Type, h events.EventHandler) {
	d.dispatcher.On(t, d.auth(h, true))
}

func (d *Dispatcher) notify(event events.Event, err error) {
	if err := d.dispatcher.Notify(event); err != nil {
		d.logger.WithError(err).Errorf("dispatch notify")
	}
}

func (d *Dispatcher) emit(event events.Event, err error) {
	if err := d.dispatcher.Emit(event); err != nil {
		d.logger.WithError(err).Errorf("dispatch emit")
	}
}

func (d *Dispatcher) reply(t events.Type, prev events.Event, payload interface{}) {
	d.notify(events.NewEvent(
		t,
		events.WithTo(prev.From),
		events.WithPayload(payload),
		events.WithPrev(&prev),
	))
}

func (d *Dispatcher) broadcast(t events.Type, prev events.Event, payload interface{}) {
	d.notify(events.NewEvent(
		t,
		events.WithTo(events.NewDestBroadcast()),
		events.WithPayload(payload),
		events.WithPrev(&prev),
	))
}

func (d *Dispatcher) auth(h events.EventHandler, required bool) events.EventHandler {
	return func(e events.Event) error {
		uid := models.NoUser
		if sid := e.GetSid(); sid != "" {
			uid, _ = d.sessionsUc.GetUid(e.Ctx, sid)
		}
		e.SetUid(uid)
		if required && uid == models.NoUser {
			d.reply(authUc.AuthRequired, e, nil)
			return usecases.ErrAuthRequired
		}
		return h(e)
	}
}
