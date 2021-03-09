package main

import (
	"context"
	"jb_chat/pkg/config"
	"jb_chat/pkg/daemon"
	"jb_chat/pkg/logger"
	"os"
)

func main() {
	appLogger := logger.DefaultLogger()
	cfg := config.MustBuild(appLogger)
	app := daemon.NewApp(cfg, appLogger)
	os.Exit(app.Run(context.Background()))
}
