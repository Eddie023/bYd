package store_test

import (
	"context"
	"testing"

	"github.com/eddie023/byd/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPosts(t *testing.T) {
	f := NewDBFixture(t)

	posts, err := f.DB.GetPosts(context.Background())
	require.NoError(t, err)

	assert.Equal(t, len(posts), 5)
}

func TestCreateNewPost(t *testing.T) {
	f := NewDBFixture(t)

	post, err := f.DB.CreateNewPost(context.Background(), store.CreatePost{
		UserID:      "user_001",
		Title:       "test title",
		Description: "test description",
		TypeID:      1,
	})
	require.NoError(t, err)

	assert.Equal(t, "user_001", post.UserID)
	assert.Equal(t, "test title", post.Title)
	assert.Equal(t, "test description", post.Description)
	assert.Equal(t, 1, post.TypeID)
}
