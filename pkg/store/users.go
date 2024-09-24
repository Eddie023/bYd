package store

import (
	"context"
	"fmt"
)

func (d *DB) GetUserByID(ctx context.Context, userID string) (User, error) {
	var user User
	err := d.pool.QueryRow(ctx, fmt.Sprintf("SELECT id FROM users WHERE id = %s", userID)).Scan(&user.ID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *DB) CreateUser(ctx context.Context, usr UserInfo) error {
	var userID string
	query := `INSERT INTO users(id, email, first_name, last_name) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING RETURNING user_id`
	err := d.pool.QueryRow(ctx, query, usr.UserId, usr.Email, usr.FirstName, usr.LastName).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}
