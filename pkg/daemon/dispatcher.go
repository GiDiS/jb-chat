package daemon

import (
	"context"
	"errors"
	"jb_chat/pkg/auth"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"jb_chat/pkg/events"
	"jb_chat/pkg/events/event"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
)

type Dispatcher struct {
	dispatcher events.Dispatcher
	c          Container
	logger     logger.Logger
	appStore   store.AppStore
}

type EventHandler func(e events.Event) ([]events.Event, error)

var googleOauthConfig *oauth2.Config

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8888/api/auth/google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func NewDispatcher(c Container) *Dispatcher {
	d := Dispatcher{
		c:          c,
		appStore:   c.Store,
		dispatcher: c.EventsDispatcher,
		logger:     c.logger,
	}
	d.init()
	return &d
}

func (d *Dispatcher) init() {
	d.dispatcher.On(event.Ping, d.onPing)
	d.dispatcher.On(event.Pong, d.onPong)
	d.dispatcher.On(event.Broadcast, d.onBroadcast)
	d.dispatcher.On(event.AuthRegister, d.onAuthRegister)
	d.dispatcher.On(event.AuthSignIn, d.onAuthSignIn)
	d.dispatcher.On(event.AuthSignOut, d.onAuthSignOut)

	d.dispatcher.On(event.ChannelsGetList, d.onChannelsGetList)
	d.dispatcher.On(event.ChannelsGetMembers, d.onChannelsGetMembers)
	d.dispatcher.On(event.UsersGetList, d.onUsersGetList)
	d.dispatcher.On(event.MessageGetList, d.onMessagesGetList)
}

func (d *Dispatcher) onPing(e events.Event) error {
	d.notify(events.NewEvent(event.Pong, events.WithTo(e.From), events.WithPrev(&e)))
	return nil
}

func (d *Dispatcher) onPong(e events.Event) error {
	d.notify(events.NewEvent(event.Ping, events.WithTo(e.From), events.WithPrev(&e)))
	return nil
}

func (d *Dispatcher) onBroadcast(e events.Event) error {
	d.notify(events.NewEvent(event.Broadcast, events.WithTo(events.NewDestBroadcast()), events.WithPrev(&e)))
	return nil
}

func (d *Dispatcher) onAuthRegister(e events.Event) error {
	payload, ok := e.Payload.(event.AuthRegisterRequest)
	if !ok {
		return errors.New("wrong req")
	}
	uid := e.GetUid()
	if uid == models.NoUser {
		return errors.New("auth required")
	}

	d.logger.Debug(payload)

	d.notify(events.NewEvent(event.AuthRegistered,
		events.WithTo(e.From),
		events.WithPayload(payload),
		events.WithPrev(&e),
	))
	d.notify(events.NewEvent(event.UsersInfo,
		events.WithTo(events.NewDestBroadcast()),
		events.WithPayload(payload),
		events.WithPrev(&e),
	))

	return nil
}

func (d *Dispatcher) onAuthSignIn(e events.Event) error {
	payload, ok := e.Payload.(event.AuthSignInRequest)
	if !ok || e.Type != event.AuthSignIn {
		return ErrInvalidRequest
	}
	var userRef *models.User
	if payload.Service == "google" {
		d.logger.Debugf("%+v", payload)
		profile, err := auth.GetProfileByAccessToken(e.Ctx, payload.AccessToken)
		if err != nil {
			return ErrAuthRequired
		}
		reqUser := profile.ToUser()
		user, err := d.upsertUser(e.Ctx, reqUser)
		if err != nil {
			return err
		}
		userRef = &user
	} else {
		d.logger.Debugf("Signin: %v", e)
	}

	if userRef != nil {
		var resp event.AuthSignInResponse
		resp.SetMe(&models.UserInfo{User: *userRef, Status: models.UserStatusUnknown})
		d.notify(events.NewEvent(
			event.AuthSignedIn,
			events.WithTo(events.NewDestBroadcast()),
			events.WithPayload(resp),
			events.WithPrev(&e),
		))
	}
	return nil
}

func (d *Dispatcher) upsertUser(ctx context.Context, newUser models.User) (models.User, error) {
	user, err := d.appStore.Users().GetByEmail(ctx, newUser.Email)

	if err != nil && err != store.ErrUserNotFound {
		return newUser, err
	}
	if err == store.ErrUserNotFound {
		uid, err := d.appStore.Users().Register(ctx, newUser)
		newUser.UserId = uid
		return newUser, err
	} else {
		_, err := d.appStore.Users().Save(ctx, user)
		return user, err
	}
}

func (d *Dispatcher) onAuthSignOut(e events.Event) error {
	d.notify(events.NewEvent(event.AuthSignedOut, events.WithTo(events.NewDestBroadcast()), events.WithPrev(&e)))
	return nil
}

func (d *Dispatcher) onChannelsGetList(e events.Event) error {
	payload := event.ChannelsListResponse{}
	channels, _ := d.appStore.Channels().Find(e.Ctx, store.ChannelsSearchCriteria{})
	payload.SetChannels(channels)
	d.notify(events.NewEvent(
		event.ChannelsList,
		events.WithTo(e.From),
		events.WithPayload(payload),
		events.WithPrev(&e)),
	)
	return nil
}

func (d *Dispatcher) onChannelsGetMembers(e events.Event) error {
	request, ok := e.Payload.(event.ChannelsMembersRequest)
	if e.Type != event.ChannelsGetMembers || !ok {
		return ErrInvalidRequest
	}
	if err := request.ChannelId.Validate(); err != nil {
		return err
	}

	members, _ := d.appStore.Members().Members(e.Ctx, request.ChannelId)
	payload := event.ChannelsMembersResponse{}
	payload.SetMembers(members)
	d.notify(events.NewEvent(
		event.ChannelsMembers,
		events.WithTo(e.From),
		events.WithPayload(payload),
		events.WithPrev(&e)),
	)
	return nil
}

func (d *Dispatcher) onUsersGetList(e events.Event) error {
	payload := event.UsersListResponse{}
	users, _ := d.appStore.Users().Find(e.Ctx, store.UserSearchCriteria{WithAvatars: true})
	infos := make([]models.UserInfo, 0, len(users))
	for _, u := range users {
		infos = append(infos, models.UserInfo{
			User: u, Status: models.UserStatusOnline,
		})
	}
	payload.SetUsers(infos)

	d.notify(events.NewEvent(
		event.UsersList,
		events.WithTo(e.From),
		events.WithPayload(payload),
		events.WithPrev(&e)),
	)
	return nil
}

func (d *Dispatcher) onMessagesGetList(e events.Event) error {
	filter := store.MessagesSearchCriteria{}
	if req, ok := e.Payload.(event.MessageListRequest); ok {
		filter.ChannelId = req.ChannelId
		if req.ParentId != nil {
			filter.ParentId = *req.ParentId
		}

	}

	messages, _ := d.appStore.Messages().Find(e.Ctx, filter)
	users := make([]models.UserInfo, 0, len(messages))
	usersMatched := make(map[models.Uid]bool)
	for _, m := range messages {
		if _, ok := usersMatched[m.UserId]; !ok {
			user, err := d.appStore.Users().GetByUid(e.Ctx, m.UserId)
			if err != nil {
				continue
			}
			users = append(users, models.UserInfo{
				User: user, Status: models.UserStatusOnline,
			})
		}
	}
	payload := event.MessageListResponse{}
	payload.SetResult(messages, users)

	d.notify(events.NewEvent(
		event.MessageList,
		events.WithTo(e.From),
		events.WithPayload(payload),
		events.WithPrev(&e)),
	)
	return nil
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
