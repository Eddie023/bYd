// package mid provides useful middleware to API handlers.
package mid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/go-chi/chi/middleware"
)

func OverrideRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		lc, ok := lambdacontext.FromContext(ctx)
		if ok {
			// override chi's request ID
			slog.Info("overriding requestID with AWS request ID", "awsRequestID", lc.AwsRequestID)
			ctx = context.WithValue(ctx, middleware.RequestIDKey, lc.AwsRequestID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
