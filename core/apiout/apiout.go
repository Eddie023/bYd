// Package apiout provides utility functions to handle API response.
package apiout

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eddie023/byd/core/logger"
	"github.com/pkg/errors"
)

func JSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("unable to marshal into JSON", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Error(ctx context.Context, w http.ResponseWriter, log *logger.Log, err error) {
	if customErr, ok := errors.Cause(err).(*CustomError); ok {
		if log != nil {
			log.ErrorWithCaller(ctx, 4, "error during request", "err", customErr.Message)
		}

		er := ErrorResponse{
			Message: customErr.Message,
			Code:    customErr.Code,
		}

		JSON(w, er, customErr.Code)
		return
	}

	if log != nil {
		log.ErrorWithCaller(ctx, 4, "unhandled error during request", "err", err)
	}
	er := ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	// if not custom error, then throw internal server error
	JSON(w, er, http.StatusInternalServerError)
}
