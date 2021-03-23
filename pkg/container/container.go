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
	authUc "github.com/GiDiS/jb-chat/pkg/usecases/auth"
	channelsUc "github.com/GiDiS/jb-chat/pkg/usecases/channels"
	dispatcherUc "github.com/GiDiS/jb-chat/pkg/usecases/dispatcher"
	messagesUc "github.com/GiDiS/jb-chat/pkg/usecases/messages"
	sessionsUc "github.com/GiDiS/jb-chat/pkg/usecases/sessions"
	systemUc "github.com/GiDiS/jb-chat/pkg/usecases/system"
	usersUc "github.com/GiDiS/jb-chat/pkg/usecases/users"
)

type Container struct {
	Config           config.Config
	Logger           logger.Logger
	Store            store.AppStore
	WsTransport      handlers_ws.WsApi
	AppDispatcher    *dispatcherUc.Dispatcher
	EventsDispatcher events.Dispatcher
	EventsResolver   events.Resolver
	AuthUsecase      authUc.Auth
	ChannelsUsecase  channelsUc.Channels
	MessagesUsecase  messagesUc.Messages
	SessionsUsecase  sessionsUc.Sessions
	SystemUsecase    systemUc.System
	UsersUsecase     usersUc.Users
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
			c.Logger.Fatalf("Failed init store", err)
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

	c.AuthUsecase = authUc.NewAuth(c.Logger, c.Store.Users())
	c.ChannelsUsecase = channelsUc.NewChannels(c.Logger, c.Store.Channels(), c.Store.Members(), c.Store.Users())
	c.MessagesUsecase = messagesUc.NewMessages(c.Logger, c.Store.Channels(), c.Store.Messages(), c.Store.Users())
	c.SessionsUsecase = sessionsUc.NewSessions(c.Logger, c.Store.Sessions(), c.Store.OnlineUsers(), c.Store.Users())
	c.SystemUsecase = systemUc.NewSystem(c.Config)
	c.UsersUsecase = usersUc.NewUsers(c.Logger, c.Store.Users())

	c.AppDispatcher = dispatcherUc.NewDispatcher(
		c.EventsDispatcher,
		c.Logger,
		c.AuthUsecase,
		c.ChannelsUsecase,
		c.MessagesUsecase,
		c.SessionsUsecase,
		c.SystemUsecase,
		c.UsersUsecase,
	)

	return &c
}
