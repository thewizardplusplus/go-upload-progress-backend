package middlewares

import (
	"net/http"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
)

func LoggingMiddleware(logger handlers.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			handler.ServeHTTP(w, r)

			logger.Printf("%s %s %s", r.Method, r.URL, time.Since(startTime))
		})
	}
}
