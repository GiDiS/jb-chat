package main

import (
	"context"
	"jb_chat/pkg/daemon"
	"jb_chat/pkg/logger"
	"os"
)

func main() {
	cfg := daemon.Config{
		PublicPort: 8888,
		DiagPort:   8889,
	}
	app := daemon.NewApp(cfg, logger.DefaultLogger())
	os.Exit(app.Run(context.Background()))
}
