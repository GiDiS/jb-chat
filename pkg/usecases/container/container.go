package container

import (
	"context"
	"jb_chat/pkg/config"
	"jb_chat/pkg/events"
	"jb_chat/pkg/handlers_ws"
	"jb_chat/pkg/logger"
	"jb_chat/pkg/store"
	"jb_chat/pkg/store/memory"
	"jb_chat/pkg/store/seed"
)

type Container struct {
	Config           config.Config
	logger           logger.Logger
	Logger           logger.Logger
	Store            store.AppStore
	WsTransport      handlers_ws.WsApi
	AppDispatcher    *Dispatcher
	EventsDispatcher events.Dispatcher
	EventsResolver   events.Resolver
}

func MustContainer(cfg config.Config, defaultLogger logger.Logger) *Container {

	c := Container{
		Logger:         defaultLogger,
		Config:         cfg,
		EventsResolver: events.DefaultResolver,
	}

	c.Store = memory.NewAppStore()

	if cfg.Seed {
		_, _ = seed.MakeSeeder(context.Background(), c.Store)
	}

	c.WsTransport = handlers_ws.NewWsTransport(c.EventsResolver, c.Logger)

	c.EventsDispatcher = events.NewDispatcher(c.Logger)
	c.EventsDispatcher.AddTransport(c.WsTransport, true)

	c.AppDispatcher = NewDispatcher(c)

	return &c
}
