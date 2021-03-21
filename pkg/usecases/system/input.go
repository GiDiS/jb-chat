package system

import "github.com/GiDiS/jb-chat/pkg/events"

type ConfigResponse struct {
	events.ResultStatus
	Config Config `json:"config"`
}

type Config struct {
	LogoUrl        string `json:"logo_url"`
	GoogleClientId string `json:"google_client_id,omitempty"`
}

func (r *ConfigResponse) SetConfig(cfg Config) {
	r.Ok = true
	r.Config = cfg
}
