package seed

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
)

type Seeder struct {
	store store.AppStore
}

func MakeSeeder(ctx context.Context, appStore store.AppStore) (*Seeder, error) {

	seeder := Seeder{
		store: appStore,
	}

	if err := seeder.fill(ctx); err != nil {
		return nil, err
	}

	return &seeder, nil
}

func (s *Seeder) fill(ctx context.Context) error {
	usersStore := s.store.Users()

	user, err := usersStore.GetByEmail(ctx, "tyrion.lannister@lannister.got")
	if err != nil && err != store.ErrUserNotFound {
		return err
	} else if user.UserId > 0 {
		return nil
	}

	chars, err := FetchCharacters()
	if err != nil {
		return err
	}

	episodes, err := FetchEpisodes()
	if err != nil {
		return err
	}

	users := BuildUsers(chars)

	for uIdx, u := range users {
		if uid, err := usersStore.Save(ctx, u); err != nil {
			return err
		} else {
			users[uIdx].UserId = uid
		}
	}

	channels, channelsMessages, channelsUsers := makeChannels(users, episodes)

	channelsStore := s.store.Channels()
	membersStore := s.store.Members()
	messagesStore := s.store.Messages()
	for _, ch := range channels {
		chUsers, _ := channelsUsers[ch.Cid]
		chMessages, _ := channelsMessages[ch.Cid]

		cid, err := channelsStore.CreatePublic(ctx, models.NoUser, ch.Title)
		if err != nil {
			return err
		}

		for uid := range chUsers {
			_ = membersStore.Join(ctx, cid, uid)
		}

		for _, msg := range chMessages {
			msg.ChannelId = cid
			_, _ = messagesStore.Create(ctx, msg)
		}
	}

	return nil
}
