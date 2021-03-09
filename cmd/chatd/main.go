package main

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/config"
	"github.com/GiDiS/jb-chat/pkg/daemon"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"os"
)

func main() {
	appLogger := logger.DefaultLogger()
	cfg := config.MustBuild(appLogger)
	app := daemon.NewApp(cfg, appLogger)
	os.Exit(app.Run(context.Background()))
}
