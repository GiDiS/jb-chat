package models

import "time"

type Subscribe struct {
	ChannelId     ChannelId  `json:"cid"`
	UserId        Uid        `json:"uid"`
	LastSeenMsgId *MessageId `json:"last_seen_mid"`
	LastSeenAt    *time.Time `json:"last_seen_at"`
	Joined        time.Time  `json:"joined"`
	Left          *time.Time `json:"left,omitempty"`
}
