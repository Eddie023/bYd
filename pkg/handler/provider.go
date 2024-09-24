package handler

import (
	"context"

	"github.com/eddie023/byd/pkg/store"
)

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

type DBProvider interface {
	GetPosts(context.Context) ([]store.Post, error)
	CreateNewPost(context.Context, store.CreatePost) (store.Post, error)
	GetUserByID(context.Context, string) (store.User, error)
}
