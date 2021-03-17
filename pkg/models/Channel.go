package models

import (
	"errors"
	"fmt"
	"time"
)

// ChannelId Positive id - direct chat with user, negative id - public channel
type ChannelId int32

func (cid ChannelId) Validate() error {
	if cid.String() == "" {
		return errors.New("empty cid")
	}
	return nil
}

const NoChannel ChannelId = 0

type ChannelType string

const (
	ChannelTypeUnknown ChannelType = ""
	ChannelTypeDirect  ChannelType = "direct"
	ChannelTypePublic  ChannelType = "public"
)

type Channel struct {
	Cid          ChannelId   `json:"cid" db:"cid"`
	Type         ChannelType `json:"type" db:"type"`
	Title        string      `json:"title"  db:"title"`
	Created      time.Time   `json:"created" db:"created"`
	OwnerUid     Uid         `json:"owner_uid"  db:"owner_uid"`
	LastMsgId    MessageId   `json:"last_msg_id" db:"last_msg_id"`
	LastMsgAt    *time.Time  `json:"last_msg_at" db:"last_msg_at"`
	MembersCount int         `json:"members_count" db:"members_count"`
}

type ChannelMembers struct {
	Members map[Uid]User `json:"Members"`
}

func (cid ChannelId) String() string {
	return string(cid)
}

func DirectTitle(uidA, uidB Uid) string {
	if uidA < uidB {
		return fmt.Sprintf("@%d:%d", uidA, uidB)
	} else {
		return fmt.Sprintf("@%d:%d", uidB, uidA)
	}
}
