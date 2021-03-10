package daemon

import (
	"context"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/config"
	"github.com/GiDiS/jb-chat/pkg/handlers_http/diag"
	"github.com/GiDiS/jb-chat/pkg/handlers_http/public"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/usecases/container"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const AppName = "jb_chat"

type App struct {
	config    config.Config
	Container *container.Container
	logger    logger.Logger
}

func NewApp(cfg config.Config, log logger.Logger) *App {
	return &App{
		config:    cfg,
		Container: container.MustContainer(cfg, log),
		logger:    log,
	}
}

func (app *App) Run(ctx context.Context) int {
	if err := app.init(ctx); err > 0 {
		return err
	}

	select {
	case <-ctx.Done():
	}
	return ErrOk
}

func (app *App) init(ctx context.Context) int {
	if err := app.initLog(ctx); err != nil {
		return ErrInitLoggerFailed
	}

	if err := app.initPublicHttpServer(ctx); err != nil {
		return ErrInitHttpServeFailed
	}

	if err := app.initDiagHttpServer(ctx); err != nil {
		return ErrInitHttpServeFailed
	}

	return ErrOk
}

func (app *App) initLog(ctx context.Context) error {
	app.logger = logger.DefaultLogger().
		WithField("app", AppName).
		WithContext(ctx)

	app.logger.Debugf("Init logger: done")
	return nil
}

func (app *App) initInterrupts(ctx context.Context) (context.Context, error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	newCtx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case sig := <-interrupt:
			app.logger.Infof("Received the OS signal %v", sig)
			cancel()
		case <-ctx.Done():
			break
		}

	}()

	app.logger.Debugf("Init interrupts: done")

	return newCtx, nil
}

func (app *App) initPublicHttpServer(ctx context.Context) error {

	publicRouter := mux.NewRouter()
	publicHandlers := public.NewRootHandlers()
	publicHandlers.RegisterHandlers(publicRouter)

	socketIoTransport := app.Container.WsTransport
	socketIoTransport.Serve(ctx)
	publicRouter.PathPrefix("/ws").HandlerFunc(socketIoTransport.ServeHTTP)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(app.config.PublicPort),
		Handler: publicRouter,
	}
	serverErrors := make(chan error)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			serverErrors <- fmt.Errorf("public http server failed: %w", err)
		}
	}()

	go func() {
		log := app.logger.WithField("facility", "http")
		select {
		case err := <-serverErrors:
			log.Errorf("Got a public http server error: %v", err)
		case <-ctx.Done():
			log.Info("Context is done")
		}

		if server == nil {
			return
		}

		timeout := 5 * time.Second
		log.Infof("Shutdown with timeout: %s", timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Server gracefully stopped")
	}()

	return nil
}

func (app *App) initDiagHttpServer(ctx context.Context) error {

	cfg := app.config
	diagRouter := mux.NewRouter()
	diagHandlers := diag.NewRootHandlers(cfg.Metrics, cfg.Pprof)
	diagHandlers.RegisterHandlers(diagRouter)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.DiagPort),
		Handler: diagRouter,
	}
	serverErrors := make(chan error)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			serverErrors <- fmt.Errorf("diag http server failed: %w", err)
		}
	}()

	go func() {
		log := app.logger.WithField("facility", "http")
		select {
		case err := <-serverErrors:
			log.Errorf("Got a diag http server error: %v", err)
		case <-ctx.Done():
			log.Info("Context is done")
		}

		if server == nil {
			return
		}

		timeout := 5 * time.Second
		log.Infof("Shutdown with timeout: %s", timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Server gracefully stopped")
	}()

	return nil
}
