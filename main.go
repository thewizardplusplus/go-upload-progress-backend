package main

import (
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

func main() {
	rand.Seed(time.Now().UnixNano())

	errorLogger := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	infoLogger := log.New(os.Stderr, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/files", handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewWritableFS("./files"),
		},
		Logger: errorLogger,
	})
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	wrappedMux := middlewares.ApplyMiddlewares(mux, []middlewares.Middleware{
		middlewares.LoggingMiddleware(infoLogger),
		middlewares.RecoveringMiddleware(errorLogger),
	})

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		errorLogger.Fatal(err)
	}
}
