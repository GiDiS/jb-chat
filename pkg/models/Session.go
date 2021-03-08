package models

import "time"

type SessionId string

type Session struct {
	UserId    Uid       `json:"uid"`
	SessionId string    `json:"sid"`
	Service   string    `json:"service"`
	AppId     string    `json:"app_id"`
	AppToken  string    `json:"app_token"`
	Token     string    `json:"token"`
	Expired   bool      `json:"expired"`
	Started   time.Time `json:"started"`
	Update    time.Time `json:"updated"`
	Expires   time.Time `json:"expires"`
	IsOnline  bool      `json:"is_online"`
}

func (sid SessionId) String() string {
	return string(sid)
}
