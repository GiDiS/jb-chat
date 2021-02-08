package models

// ChannelId Positive id - direct chat with user, negative id - public channel
type ChannelId int32

const NoChannel ChannelId = 0

type ChannelType uint8

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeDirect
	ChannelTypePublic
)

type Channel struct {
	Cid          ChannelId   `json:"cid"`
	Title        string      `json:"title"`
	LastMsg      MessageId   `json:"last_msg"`
	MembersCount uint32      `json:"members_count"`
	Type         ChannelType `json:"type"`
}

type ChannelMembers struct {
	Members map[Uid]User `json:"Members"`
}
