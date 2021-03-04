package seed

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAnswer(t *testing.T) {

	answer, err := GetAnswer(context.Background(), "And about established lying? rules the we've")
	assert.NoError(t, err)
	assert.NotEmpty(t, answer)
}
