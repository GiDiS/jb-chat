package container

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/config"
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/handlers_ws"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/store/memory"
	"github.com/GiDiS/jb-chat/pkg/store/postgres"
	"github.com/GiDiS/jb-chat/pkg/store/postgres/migration"
	"github.com/GiDiS/jb-chat/pkg/store/seed"
)

type Container struct {
	Config           config.Config
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

	if c.Config.GetStore() == config.StorePostgres {
		dbConfig := postgres.MustGetConfig(c.Logger)
		db := postgres.ConnectToDB("postgres", dbConfig.Dsn(), c.Logger)
		migration.MustMigrate(db)
		appStore, err := postgres.NewAppStore(db)
		if err != nil {
			log.Fatalf("Failed init store", err)
		}

		c.Store = appStore
	} else {
		c.Store = memory.NewAppStore()
	}

	if cfg.Seed {
		c.Logger.Debug("Start seeding")
		_, err := seed.MakeSeeder(context.Background(), c.Store)
		if err != nil {
			c.Logger.Fatalf("Seeding failed: %v", err)
		}
		c.Logger.Debug("Finish seeding")
	}

	c.WsTransport = handlers_ws.NewWsTransport(c.EventsResolver, c.Logger)

	c.EventsDispatcher = events.NewDispatcher(c.Logger)
	c.EventsDispatcher.AddTransport(c.WsTransport, true)

	c.AppDispatcher = NewDispatcher(c)

	return &c
}
