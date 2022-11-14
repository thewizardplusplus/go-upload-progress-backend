package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers/middlewares"
	"github.com/thewizardplusplus/go-upload-progress/gateways/writablefs"
	"github.com/thewizardplusplus/go-upload-progress/usecases"
)

//go:generate swag init --dir ../../ --generalInfo ./cmd/go-upload-progress/main.go --output ../../docs/ --outputTypes yaml --propertyStrategy pascalcase

// @title go-upload-progress API
// @version 1.0.0
// @license.name MIT
// @host localhost:8080
// @basePath /api/v1

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

	server := http.Server{
		Addr: serverAddress,
		Handler: middlewares.ApplyMiddlewares(mux, []middlewares.Middleware{
			middlewares.LoggingMiddleware(infoLogger),
			middlewares.RecoveringMiddleware(errorLogger),
		}),
	}

	// https://pkg.go.dev/net/http@go1.19.0#Server.Shutdown
	//
	// BSD 3-Clause "New" or "Revised" License
	// Copyright (C) 2009 The Go Authors
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
