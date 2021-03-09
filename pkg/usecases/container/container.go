package container

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/config"
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/handlers_ws"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/store/memory"
	"github.com/GiDiS/jb-chat/pkg/store/seed"
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
