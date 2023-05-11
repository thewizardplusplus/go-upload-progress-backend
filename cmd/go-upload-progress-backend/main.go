package main

// nolint: lll
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/thewizardplusplus/go-upload-progress-backend/gateways/handlers"
	"github.com/thewizardplusplus/go-upload-progress-backend/gateways/handlers/middlewares"
	writablefs "github.com/thewizardplusplus/go-upload-progress-backend/gateways/writable-fs"
	"github.com/thewizardplusplus/go-upload-progress-backend/usecases"
	"github.com/thewizardplusplus/go-upload-progress-backend/usecases/generators"
)

// @title go-upload-progress-backend API
// @version 1.4.0
// @license.name MIT
// @host localhost:8080
// @basePath /api/v1

const (
	uploadedFileRoute = "/files/"
	loggerFlags       = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix
)

func main() {
	infoLogger := makeLogger("INFO")
	errorLogger := makeLogger("ERROR")

	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	staticFileDir := getEnv("STATIC_FILE_DIR", "./static")
	uploadedFileDir := getEnv("UPLOADED_FILE_DIR", "./files")

	randomSuffixByteCountAsString := getEnv("RANDOM_SUFFIX_BYTE_COUNT", "4")
	randomSuffixByteCount, err := strconv.Atoi(randomSuffixByteCountAsString)
	if err != nil {
		errorLogger.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/files", handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewDirFS(uploadedFileDir),
			FilenameGenerator: generators.FilenameGenerator{
				RandomSuffixByteCount: randomSuffixByteCount,
			},
		},
		Logger: errorLogger,
	})
	mux.Handle(
		uploadedFileRoute,
		http.StripPrefix(uploadedFileRoute, makeFileServer(uploadedFileDir)),
	)
	mux.Handle("/", makeFileServer(staticFileDir))

	server := http.Server{
		Addr: serverAddress,
		Handler: middlewares.ApplyMiddlewares(mux, []middlewares.Middleware{
			middlewares.LoggingMiddleware(infoLogger),
			middlewares.RecoveringMiddleware(errorLogger),
		}),
	}

	// https://pkg.go.dev/net/http@go1.19.0#example-Server.Shutdown
	//
	// # License
	//
	//	BSD 3-Clause "New" or "Revised" License
	//	Copyright (C) 2009 The Go Authors
	waitingToShutdown := make(chan struct{})
	go func() {
		defer close(waitingToShutdown)

		waitingToInterrupt := make(chan os.Signal, 1)
		signal.Notify(waitingToInterrupt, os.Interrupt)
		<-waitingToInterrupt

		if err := server.Shutdown(context.Background()); err != nil {
			errorLogger.Print(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		errorLogger.Fatal(err)
	}
	<-waitingToShutdown
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
