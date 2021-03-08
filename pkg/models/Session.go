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
	Expires   time.Time `json:"expires"`
}

func (sid SessionId) String() string {
	return string(sid)
}
