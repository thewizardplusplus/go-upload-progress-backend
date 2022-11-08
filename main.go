package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/thewizardplusplus/go-upload-progress-backend/gateways/handlers"
	writablefs "github.com/thewizardplusplus/go-upload-progress-backend/gateways/writable-fs"
	"github.com/thewizardplusplus/go-upload-progress-backend/usecases"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fileHandler := handlers.FileHandler{
		FileUsecase: usecases.FileUsecase{
			WritableFS: writablefs.NewWritableFS("./files"),
		},
	}

	http.Handle("/api/v1/files", fileHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	http.Handle("/", http.FileServer(http.Dir("./public")))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
