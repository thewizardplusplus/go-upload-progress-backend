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

	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	publicFileDir := getEnv("PUBLIC_FILE_DIR", "./public")
	uploadedFileDir := getEnv("UPLOADED_FILE_DIR", "./files")

	infoLogger := makeLogger("INFO")
	errorLogger := makeLogger("ERROR")

	mux := http.NewServeMux()
	mux.Handle("/api/v1/files", handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewWritableFS(uploadedFileDir),
		},
		Logger: errorLogger,
	})
	mux.Handle(
		uploadedFileRoute,
		http.StripPrefix(uploadedFileRoute, makeFileServer(uploadedFileDir)),
	)
	mux.Handle("/", makeFileServer(publicFileDir))

	wrappedMux := middlewares.ApplyMiddlewares(mux, []middlewares.Middleware{
		middlewares.LoggingMiddleware(infoLogger),
		middlewares.RecoveringMiddleware(errorLogger),
	})

	if err := http.ListenAndServe(serverAddress, wrappedMux); err != nil {
		errorLogger.Fatal(err)
	}
}

func getEnv(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = defaultValue
	}

	return value
}

func makeLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s] ", prefix), loggerFlags)
}

func makeFileServer(baseDir string) http.Handler {
	return http.FileServer(http.Dir(baseDir))
}
