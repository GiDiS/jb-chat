package system

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/config"
)

type System interface {
	GetConfig(ctx context.Context) (ConfigResponse, error)
}

type systemImpl struct {
	cfg config.Config
}

func NewSystem(cfg config.Config) *systemImpl {
	return &systemImpl{cfg: cfg}
}
