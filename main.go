package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers/middlewares"
	"github.com/thewizardplusplus/go-upload-progress/gateways/writablefs"
	"github.com/thewizardplusplus/go-upload-progress/usecases"
)

const (
	uploadedFileRoute = "/files/"
	loggerFlags       = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix
)

func main() {
	rand.Seed(time.Now().UnixNano())

	infoLogger := makeLogger("INFO")
	errorLogger := makeLogger("ERROR")

	mux := http.NewServeMux()
	mux.Handle("/api/v1/files", handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewWritableFS("./files"),
		},
		Logger: errorLogger,
	})
	mux.Handle(uploadedFileRoute, http.StripPrefix(uploadedFileRoute, makeFileServer("./files")))
	mux.Handle("/", makeFileServer("./public"))

	wrappedMux := middlewares.ApplyMiddlewares(mux, []middlewares.Middleware{
		middlewares.LoggingMiddleware(infoLogger),
		middlewares.RecoveringMiddleware(errorLogger),
	})

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		errorLogger.Fatal(err)
	}
}

func makeLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s] ", prefix), loggerFlags)
}

func makeFileServer(baseDir string) http.Handler {
	return http.FileServer(http.Dir(baseDir))
}
