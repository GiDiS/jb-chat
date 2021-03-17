package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type MessageId uint32

const NoMessage MessageId = 0

type Message struct {
	ChannelId ChannelId    `json:"cid" db:"cid"`
	MsgId     MessageId    `json:"mid" db:"mid"`
	UserId    Uid          `json:"uid" db:"uid"`
	ParentId  MessageId    `json:"pid,omitempty" db:"pid"`
	Created   time.Time    `json:"created" db:"created"`
	Deleted   *time.Time   `json:"deleted,omitempty" db:"deleted"`
	Body      string       `json:"body" db:"body"`
	IsThread  bool         `json:"is_thread" db:"is_thread"`
	Likes     MessageLikes `json:"likes" db:"likes"`
}

type MessageLikes []Uid

func (l *MessageLikes) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*l = make(MessageLikes, 0)
	if likes, ok := value.(string); ok {
		likes = strings.Trim(likes, "[]")
		for _, s := range strings.Split(likes, ",") {
			uid, _ := strconv.Atoi(s)
			*l = append(*l, Uid(uid))
		}
		return nil
	}
	return errors.New("invalid likes")
}
