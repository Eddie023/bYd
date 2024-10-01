package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/eddie023/byd/core/logger"
	"github.com/eddie023/byd/pkg/auth"
	"github.com/eddie023/byd/pkg/config"
	"github.com/eddie023/byd/pkg/handler"
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

type Lambda struct {
	Server http.Handler
}

func (h *Lambda) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adapter := handlerfunc.New(h.Server.ServeHTTP)

	return adapter.ProxyWithContext(ctx, req)
}

func buildHandler() (*Lambda, error) {
	ctx := context.Background()
	log := logger.New(os.Stdout, "lambda-rest-api")

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	log.Debug(ctx, "using config", "config", cfg)
	db, err := store.NewDB(ctx, cfg.Db.ConnectionURI)
	if err != nil {
		return nil, err
	}

	h, err := handler.NewAPIHandler(db, log, &auth.LambdaAuthenticator{})
	if err != nil {
		return nil, err
	}

	lambda := Lambda{
		Server: h,
	}

	return &lambda, nil
}
