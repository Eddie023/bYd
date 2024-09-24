package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUsers(t *testing.T) {
	f := NewDBFixture(t)

	user, err := f.DB.GetUserByID(context.Background(), "user_001")
	require.NoError(t, err)

	assert.Equal(t, "user_001", user.ID)
}
