package models

import "time"

type MessageId uint32

const NoMessage MessageId = 0

type Message struct {
	ChannelId ChannelId `json:"cid"`
	MsgId     MessageId `json:"mid"`
	ParentId  MessageId `json:"pid,omitempty"`
	Created   time.Time `json:"created"`
	Body      string    `json:"body"`
	IsThread  bool      `json:"is_thread"`
}
