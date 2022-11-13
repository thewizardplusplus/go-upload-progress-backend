package middlewares

import (
	"net/http"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
)

const (
	defaultStatusCode = http.StatusOK
)

type responseWriterWrapper struct {
	http.ResponseWriter

	statusCode       int
	writtenByteCount int
	isWritten        bool
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter:   w,
		statusCode:       defaultStatusCode,
		writtenByteCount: 0,
		isWritten:        false,
	}
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if w.isWritten {
		return
	}

	w.statusCode = statusCode
	w.isWritten = true

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	w.isWritten = true

	writtenByteCount, err := w.ResponseWriter.Write(data)
	w.writtenByteCount += writtenByteCount

	return writtenByteCount, err
}

func LoggingMiddleware(logger handlers.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			wrapper := newResponseWriterWrapper(w)
			handler.ServeHTTP(wrapper, r)

			logger.Printf(
				"%s %s %d %dB %s",
				r.Method,
				r.URL,
				wrapper.statusCode,
				wrapper.writtenByteCount,
				time.Since(startTime),
			)
		})
	}
}
