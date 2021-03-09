package seed

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const CharactersUrl = "https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/characters.json"
const EpisodesUrl = "https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/episodes.json"
const WordsUrl = "https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/script-bag-of-words.json"
const LordVarys = "Lord Varys"

type CharactersList struct {
	Characters []Character `json:"characters"`
}

type Character struct {
	CharacterName       string     `json:"characterName"`
	HouseNames          HouseNames `json:"houseName"`
	CharacterImageThumb string     `json:"characterImageThumb"`
	CharacterImageFull  string     `json:"characterImageFull"`
	CharacterLink       string     `json:"characterLink"`
	Nickname            string     `json:"nickname"`
	Killed              []string   `json:"killed"`
	ServedBy            []string   `json:"servedBy"`
	ParentOf            []string   `json:"parentOf"`
	Siblings            []string   `json:"siblings"`
	KilledBy            []string   `json:"killedBy"`
}

type EpisodesList struct {
	Episodes []Episode
}

type HouseNames []string

type Episode struct {
	EpisodeAlt   string     `json:"episodeAlt"`
	SeasonNum    int        `json:"seasonNum"`
	EpisodeNum   int        `json:"episodeNum"`
	EpisodeTitle string     `json:"episodeTitle"`
	Items        []TextItem `json:"text"`
}

type TextItem struct {
	Text string `json:"text"`
	Name string `json:"name"`
}

func (h *HouseNames) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	if string(bytes[0]) == "\"" {
		var name string
		if err := json.Unmarshal(bytes, &name); err != nil {
			return err
		}
		*h = []string{name}
	} else if string(bytes[0]) == "[" {
		var names []string
		if err := json.Unmarshal(bytes, &names); err != nil {
			return err
		}
		*h = names
	}
	return nil
}

func FetchCharacters() ([]Character, error) {
	resp, err := http.Get(CharactersUrl)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return nil, errors.New("empty resp")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var chars CharactersList
	err = json.Unmarshal(body, &chars)
	if err != nil {
		fmt.Println(body[21275:])
		return nil, err
	}
	return chars.Characters, nil
}

func FetchEpisodes() ([]Episode, error) {
	resp, err := http.Get(WordsUrl)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return nil, errors.New("empty resp")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var episodes []Episode
	err = json.Unmarshal(body, &episodes)
	if err != nil {
		return nil, err
	}
	return episodes, nil
}

func BuildUsers(characters []Character) []models.User {
	users := make([]models.User, 0, len(characters))

	for idx, ch := range characters {
		user := makeUser(ch)
		user.UserId = models.Uid(idx + 1)
		users = append(users, user)
	}
	return users
}

func makeUser(ch Character) models.User {
	domain := "other.got"
	if len(ch.HouseNames) > 0 {
		domain = ch.HouseNames[0] + ".got"
	}
	re := regexp.MustCompile(`\W`)
	login := re.ReplaceAllString(ch.CharacterName, ".")
	nickname := regexp.MustCompile(`\W`).ReplaceAllString(ch.CharacterName, "")

	return models.User{
		Nickname:  nickname,
		Title:     ch.CharacterName,
		AvatarUrl: ch.CharacterImageFull,
		Email:     strings.ToLower(login + "@" + domain),
	}
}

func makeChannels(users []models.User, episodes []Episode) (
	[]models.Channel, map[models.ChannelId][]models.Message, map[models.ChannelId]map[models.Uid]models.MessageId,
) {
	userMap := mapUsersByName(users)

	channels := make([]models.Channel, 0, len(episodes))
	channelsMessages := make(map[models.ChannelId][]models.Message, len(episodes))
	channelsUsers := make(map[models.ChannelId]map[models.Uid]models.MessageId, len(episodes))
	for epIdx, ep := range episodes {
		ch := models.Channel{
			Cid:   models.ChannelId(epIdx + 1),
			Title: fmt.Sprintf("S%dE%02d: %s", ep.SeasonNum, ep.EpisodeNum, ep.EpisodeTitle),
			Type:  models.ChannelTypePublic,
		}
		timeOffset := time.Now().Add(-time.Duration(len(ep.Items)) * time.Minute)
		channelUsers := make(map[models.Uid]models.MessageId)
		channelMessages := make([]models.Message, 0, len(ep.Items))
		msgSeq := 1
		for _, txt := range ep.Items {
			msgId := models.MessageId(msgSeq)
			ch.LastMsg = msgId
			uid := models.NoUser
			msgSeq++

			if user, ok := userMap[txt.Name]; ok {
				uid = user.UserId
				channelUsers[uid] = msgId
			}

			channelMessages = append(channelMessages, models.Message{
				ChannelId: ch.Cid,
				MsgId:     msgId,
				UserId:    uid,
				ParentId:  models.NoMessage,
				Created:   timeOffset.Add(time.Minute * time.Duration(msgSeq)),
				Body:      txt.Text,
				IsThread:  false,
			})

		}
		if l := len(channelMessages) - 1; l >= 0 {
			last := channelMessages[l]
			ch.LastMsgAt = last.Created
			ch.LastMsg = last.MsgId
		}
		ch.MembersCount = len(channelUsers)
		channels = append(channels, ch)
		channelsMessages[ch.Cid] = channelMessages
		channelsUsers[ch.Cid] = channelUsers
	}

	return channels, channelsMessages, channelsUsers
}

func mapUsersByName(users []models.User) map[string]models.User {
	usersMap := make(map[string]models.User, len(users))
	for _, user := range users {
		usersMap[user.Title] = user
	}
	return usersMap
}
