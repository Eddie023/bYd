package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/eddie023/byd/core/logger"
	"github.com/eddie023/byd/internal/build"
	"github.com/eddie023/byd/pkg/auth"
	"github.com/eddie023/byd/pkg/config"
	"github.com/eddie023/byd/pkg/handler"
	"github.com/eddie023/byd/pkg/store"
)

func main() {
	ctx := context.Background()

	log := logger.New(os.Stderr, "api-server", logger.WithColor())
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Log) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build.Build)

	cfg, err := config.New()
	if err != nil {
		return err
	}

	log.Debug(ctx, "using config", "config", cfg)
	db, err := store.NewDB(ctx, cfg.Db.ConnectionURI)
	if err != nil {
		return fmt.Errorf("connecting db: %w", err)
	}

	h, err := handler.NewAPIHandler(db, log, &auth.LocalAuth{
		UserID: "user_001",
	})
	if err != nil {
		return fmt.Errorf("failed setting up handler: %w", err)
	}

	api := &http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: h,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "server listening on", "port", cfg.Web.APIHost, "host", cfg.Web.APIHost)

		serverErrors <- api.ListenAndServe()
	}()

	// Listen for syscall signals for process to interrput/quit
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server errror: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Error(ctx, "graceful shutdown failed", "deadline exceeded", true)
			}
		}()

		// Trigger graceful shutdown
		log.Warn(ctx, "gracefully shutting down server", "deadline exceeded", false)
		if err := api.Shutdown(shutdownCtx); err != nil {
			if err = api.Close(); err != nil {
				return err
			}
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
