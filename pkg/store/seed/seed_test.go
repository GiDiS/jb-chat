package seed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchCharacters(t *testing.T) {
	got, err := FetchCharacters()
	assert.NoError(t, err)
	assert.NotZero(t, len(got))
}

func TestFetchEpisodes(t *testing.T) {
	got, err := FetchEpisodes()
	assert.NoError(t, err)
	assert.NotZero(t, len(got))
}

func Test_makeChannels(t *testing.T) {
	chars, err := FetchCharacters()
	assert.NoError(t, err)

	episodes, err := FetchEpisodes()
	assert.NoError(t, err)

	users := BuildUsers(chars)

	channels, chMessages, chUsers := makeChannels(users, episodes)
	assert.NotZero(t, len(channels))
	assert.NotZero(t, len(chMessages))
	assert.NotZero(t, len(chUsers))
}
