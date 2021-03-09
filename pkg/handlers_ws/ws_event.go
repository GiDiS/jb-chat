package handlers_ws

import (
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/models"
	"time"
)

const (
	WsConnected    events.Type = "ws.connected"
	WsDisconnected events.Type = "ws.disconnected"
)

type SysClientResponse struct {
	Id             string     `json:"id"`
	RemoteAddr     string     `json:"remote_addr"`
	LocalAddr      string     `json:"local_addr"`
	LocalPort      int        `json:"local_port"`
	Online         bool       `json:"online"`
	ConnectedAt    time.Time  `json:"connected"`
	DisconnectedAt *time.Time `json:"disconnected"`
	Uid            models.Uid `json:"uid"`
}
