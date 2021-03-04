package daemon

import (
	"context"
	"jb_chat/pkg/events"
	"jb_chat/pkg/handlers_ws"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/store"
	"jb_chat/pkg/store/memory"
	"jb_chat/pkg/store/seed"
)

type Container struct {
	Config           Config
	logger           logger.Logger
	Store            store.AppStore
	WsTransport      handlers_ws.WsApi
	AppDispatcher    *Dispatcher
	EventsDispatcher events.Dispatcher
	EventsResolver   events.Resolver
}

func MustContainer(cfg Config, defaultLogger logger.Logger) *Container {

	c := Container{
		logger:         defaultLogger,
		Config:         cfg,
		EventsResolver: events.DefaultResolver,
	}

	c.Store = memory.NewAppStore()

	_, _ = seed.MakeSeeder(context.Background(), c.Store)

	c.WsTransport = handlers_ws.NewWsTransport(c.EventsResolver, c.logger)

	c.EventsDispatcher = events.NewDispatcher(c.logger)
	c.EventsDispatcher.AddTransport(c.WsTransport, true)

	c.AppDispatcher = NewDispatcher(c)

	return &c
}
