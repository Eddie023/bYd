package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/eddie023/byd/core/apiout"
	"github.com/eddie023/byd/core/logger"
	"github.com/eddie023/byd/pkg/auth"
	"github.com/eddie023/byd/pkg/mid"
	"github.com/eddie023/byd/pkg/types"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/rs/cors"

	httpmid "github.com/oapi-codegen/nethttp-middleware"
)

type APIHandler struct {
	db         DBProvider
	authorizer auth.Authenticator
	log        *logger.Log
}

func NewAPIHandler(db DBProvider, log *logger.Log, authorizer auth.Authenticator) (http.Handler, error) {
	swagger, err := types.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed setting up swagger spec %w", err)
	}

	swagger.Servers = nil
	a := APIHandler{
		db:         db,
		log:        log,
		authorizer: authorizer,
	}

	r := chi.NewRouter()
	r.Use(mid.OverrideRequestID)

	chiLogger := &httplog.Logger{
		Logger: slog.New(log.Handler),
		Options: httplog.Options{
			LogLevel:        slog.LevelDebug,
			QuietDownPeriod: 10 * time.Second,
		},
	}

	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(chiLogger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(crs.Handler)
	r.Use(middleware.Heartbeat("/alive"))

	r.Group(func(r chi.Router) {
		r.Use(httpmid.OapiRequestValidatorWithOptions(swagger, &httpmid.Options{
			// pass the custom error handler such that we return JSON response instead of default text
			// for validation error
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				apiout.JSON(w, apiout.ErrorResponse{
					Code:    statusCode,
					Message: message,
				}, http.StatusBadRequest)
			},
		}))
		r.Use(auth.Middleware(a.authorizer, a.db, a.log))
		types.HandlerWithOptions(&a, types.ChiServerOptions{
			BaseRouter: r,
		})
	})

	return r, nil
}
