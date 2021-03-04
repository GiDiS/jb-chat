package models

import (
	"errors"
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

type ChannelType uint8

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeDirect
	ChannelTypePublic
)

type Channel struct {
	Cid          ChannelId   `json:"cid"`
	OwnerId      Uid         `json:"owner_id"`
	Title        string      `json:"title"`
	LastMsg      MessageId   `json:"last_msg_id"`
	LastMsgAt    time.Time   `json:"last_msg_at"`
	MembersCount int         `json:"members_count"`
	Type         ChannelType `json:"type"`
}

type ChannelMembers struct {
	Members map[Uid]User `json:"Members"`
}

func (cid ChannelId) String() string {
	return string(cid)
}
