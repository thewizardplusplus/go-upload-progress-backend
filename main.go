package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
	"github.com/thewizardplusplus/go-upload-progress/gateways/writablefs"
	"github.com/thewizardplusplus/go-upload-progress/usecases"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/api/v1/files", func(w http.ResponseWriter, r *http.Request) {
		fileHandler := handlers.FileHandler{
			FileUsecase: usecases.FileUsecase{
				WritableFS: writablefs.NewWritableFS("./files"),
			},
		}

		switch r.Method {
		case http.MethodGet:
			fileHandler.GetFiles(w, r)
		case http.MethodPost:
			fileHandler.SaveFile(w, r)
		case http.MethodDelete:
			if filename := r.FormValue("filename"); filename != "" {
				fileHandler.DeleteFile(w, r)
			} else {
				fileHandler.DeleteFiles(w, r)
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	http.Handle("/", http.FileServer(http.Dir("./public")))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
