package models

import "time"

type SessionId string

type Session struct {
	SessionId string    `json:"sid" db:"sid"`
	UserId    Uid       `json:"uid" db:"uid"`
	Service   string    `json:"service" db:"service"`
	AppId     string    `json:"app_id" db:"app_id"`
	AppToken  string    `json:"app_token" db:"app_token"`
	Token     string    `json:"token" db:"token"`
	Expired   bool      `json:"expired" db:"expired"`
	IsOnline  bool      `json:"is_online" db:"is_online"`
	Started   time.Time `json:"started" db:"started"`
	Updated   time.Time `json:"updated" db:"updated"`
	Expires   time.Time `json:"expires" db:"expires"`
}

func (sid SessionId) String() string {
	return string(sid)
}
