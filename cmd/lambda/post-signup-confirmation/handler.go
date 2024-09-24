package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/eddie023/byd/pkg/config"
	"github.com/eddie023/byd/pkg/store"
)

var l *Lambda

func init() {
	var err error
	l, err = buildHandler()
	if err != nil {
		panic(err)
	}
}

func main() {
	awslambda.Start(l.Handler)
}

type DBProvider interface {
	CreateUser(ctx context.Context, usr store.UserInfo) error
}

type Lambda struct {
	db DBProvider
}

func (h *Lambda) Handler(ctx context.Context, req events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	slog.Info("post confirmation event", "req", req)

	sub, ok := req.Request.UserAttributes["sub"]
	if !ok {
		return req, errors.New("could not find usb field")
	}

	email, ok := req.Request.UserAttributes["email"]
	if !ok {
		return req, errors.New("could not find email field")
	}

	givenName, ok := req.Request.UserAttributes["given_name"]
	if !ok {
		return req, errors.New("could not find 'given_name' field")
	}

	familyName, ok := req.Request.UserAttributes["family_name"]
	if !ok {
		return req, errors.New("could not find 'family_name' field")
	}

	err := h.db.CreateUser(ctx, store.UserInfo{
		UserId:    sub,
		Email:     email,
		FirstName: givenName,
		LastName:  familyName,
	})
	if err != nil {
		slog.Error("unable to store in db", "err", err)
		return req, err
	}

	slog.Info("successfully created new user")
	return req, nil
}

func buildHandler() (*Lambda, error) {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	slog.Info("using config", "uri", cfg.Db.ConnectionURI)
	db, err := store.NewDB(ctx, cfg.Db.ConnectionURI)
	if err != nil {
		return nil, err
	}

	lambda := Lambda{
		db: db,
	}

	return &lambda, nil
}
