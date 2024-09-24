package auth

import "net/http"

// LocalAuth
type LocalAuth struct {
	UserID string
	Email  string
}

func NewLocalAuth(userID, email string) *LocalAuth {
	return &LocalAuth{
		UserID: userID,
		Email:  email,
	}
}

func (l *LocalAuth) Authenticate(r *http.Request) (Claims, error) {
	return Claims{
		Sub:   l.UserID,
		Email: l.Email,
	}, nil
}
