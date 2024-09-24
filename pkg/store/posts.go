package store

import "context"

func (d *DB) GetPosts(ctx context.Context) ([]Post, error) {
	var posts []Post

	query := `SELECT id, user_id, title, is_anon, description, type_id FROM post`
	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.IsAnon, &p.Description, &p.TypeID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

type CreatePost struct {
	UserID      string
	Title       string
	IsAnon      bool
	Description string
	TypeID      int
}

func (d *DB) CreateNewPost(ctx context.Context, post CreatePost) (Post, error) {
	var p Post

	query := `INSERT INTO post (user_id, title, is_anon, description, type_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	RETURNING (id, user_id, title, is_anon, description, type_id, created_at, updated_at)`
	err := d.pool.QueryRow(ctx, query, post.UserID, post.Title, post.IsAnon, post.Description, post.TypeID).Scan(&p)
	if err != nil {
		return Post{}, err
	}

	return p, nil
}
