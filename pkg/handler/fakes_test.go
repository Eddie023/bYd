package handler_test

import (
	"context"

	"github.com/eddie023/byd/pkg/store"
)

var TestUser = store.User{
	ID:    "user_001",
	Email: "user_001@gmail.com",
}

type FakeDBProvider struct {
	err           error
	createNewPost *store.CreatePost
}

func (f *FakeDBProvider) GetPosts(context.Context) ([]store.Post, error) {
	return nil, nil
}

func (f *FakeDBProvider) CreateNewPost(ctx context.Context, p store.CreatePost) (store.Post, error) {
	if f.createNewPost == nil {
		return store.Post{}, f.err
	}

	return store.Post{
		ID:          1,
		TypeID:      1,
		UserID:      TestUser.ID,
		Title:       f.createNewPost.Title,
		Description: f.createNewPost.Description,
		IsAnon:      f.createNewPost.IsAnon,
	}, f.err
}

func (f *FakeDBProvider) GetUserByID(ctx context.Context, usr string) (store.User, error) {
	if usr == TestUser.ID {
		return TestUser, nil
	}

	return store.User{}, nil
}
