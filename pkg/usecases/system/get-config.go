package system

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/events"
)

func (u *systemImpl) GetConfig(ctx context.Context) (ConfigResponse, error) {
	return ConfigResponse{
		Config: Config{
			LogoUrl:        "/ui/logo.png",
			GoogleClientId: u.cfg.GoogleAuth.ClientID,
		},
		ResultStatus: events.ResultStatus{Ok: true},
	}, nil

}
