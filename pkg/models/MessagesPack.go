package models

// MessagePack - bulk pack of chanel messages with embedded users info
type MessagePack struct {
	ChannelId ChannelId        `json:"cid"`
	Messages  []Message        `json:"messages"`
	Users     map[Uid]UserInfo `json:"users"`
}
