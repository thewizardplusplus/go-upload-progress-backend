package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
	"github.com/thewizardplusplus/go-upload-progress/gateways/writablefs"
	"github.com/thewizardplusplus/go-upload-progress/usecases"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	fileHandler := handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewWritableFS("./files"),
		},
		Logger: logger,
	}

	http.Handle("/api/v1/files", fileHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	http.Handle("/", http.FileServer(http.Dir("./public")))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal(err)
	}
}
