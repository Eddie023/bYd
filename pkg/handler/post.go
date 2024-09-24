package handler

import (
	"encoding/json"
	"net/http"

	"github.com/eddie023/byd/core/apiout"
	"github.com/eddie023/byd/pkg/auth"
	"github.com/eddie023/byd/pkg/store"
	"github.com/eddie023/byd/pkg/types"
)

// (GET /v1/post
func (a *APIHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := a.db.GetPosts(ctx)
	if err != nil {
		apiout.Error(ctx, w, a.log, err)
		return
	}

	var postsJSON []types.Post
	for _, p := range posts {
		postsJSON = append(postsJSON, types.Post{
			Id:          p.ID,
			Description: p.Description,
			UserId:      p.UserID,
			IsAnon:      p.IsAnon,
			Title:       p.Title,
		})
	}

	apiout.JSON(w, postsJSON, http.StatusOK)
}

// (POST /v1/posts
func (a *APIHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := types.CreatePost{}

	userId, err := auth.GetUserID(ctx)
	if err != nil {
		apiout.Error(ctx, w, a.log, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiout.Error(ctx, w, a.log, err)
		return
	}

	post, err := a.db.CreateNewPost(ctx, store.CreatePost{
		UserID:      userId,
		Title:       request.Title,
		Description: request.Description,
	})
	if err != nil {
		apiout.Error(ctx, w, a.log, err)
		return
	}

	response := types.Post{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		UserId:      post.UserID,
		IsAnon:      post.IsAnon,
		CreatedAt:   post.CreatedAt,
		Type:        post.TypeID,
	}

	apiout.JSON(w, response, http.StatusOK)
}
