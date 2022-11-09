package middlewares

import (
	"net/http"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
)

func RecoveringMiddleware(logger handlers.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Print(err)
				}
			}()

			handler.ServeHTTP(w, r)
		})
	}
}
