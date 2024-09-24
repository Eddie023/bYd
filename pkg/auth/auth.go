package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eddie023/byd/core/apiout"
	"github.com/eddie023/byd/core/logger"
	"github.com/eddie023/byd/pkg/store"
)

type ctxKey int

const (
	userIDKey ctxKey = iota + 1
)

func setUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (string, error) {
	v, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", errors.New("user id not found in context")
	}

	return v, nil
}

type Claims struct {
	Sub   string
	Email string
}

type Authenticator interface {
	Authenticate(r *http.Request) (Claims, error)
}

type AuthDBProvider interface {
	GetUserByID(ctx context.Context, userID string) (store.User, error)
}

func Middleware(authenticator Authenticator, db AuthDBProvider, log *logger.Log) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			claims, err := authenticator.Authenticate(r)
			if err != nil {
				slog.Error("middleware", "err", err)
				apiout.Error(ctx, w, log, errors.New("authenticator error"))
				return
			}

			// lookup user information in the db from claims
			user, err := db.GetUserByID(ctx, claims.Sub)
			if err != nil {
				apiout.Error(ctx, w, log, err)
				return
			}

			ctx = setUserID(ctx, user.ID)
			slog.Debug("user is authenticated", "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
