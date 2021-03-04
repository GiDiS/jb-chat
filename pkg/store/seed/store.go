package seed

import (
	"context"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
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
	chars, err := FetchCharacters()
	if err != nil {
		return err
	}

	episodes, err := FetchEpisodes()
	if err != nil {
		return err
	}

	users := BuildUsers(chars)

	usersStore := s.store.Users()
	for _, u := range users {
		if _, err := usersStore.Save(ctx, u); err != nil {
			return err
		}
	}

	channels, channelsMessages, channelsUsers := makeChannels(users, episodes)

	channelsStore := s.store.Channels()
	membersStore := s.store.Members()
	messagesStore := s.store.Messages()
	for _, ch := range channels {
		chUsers, _ := channelsUsers[ch.Cid]
		chMessages, _ := channelsMessages[ch.Cid]

		_, _ = channelsStore.CreatePublic(ctx, models.NoUser, ch.Title)

		for uid := range chUsers {
			_ = membersStore.Join(ctx, ch.Cid, uid)
		}

		for _, msg := range chMessages {
			_, _ = messagesStore.Create(ctx, msg)
		}
	}

	return nil
}
