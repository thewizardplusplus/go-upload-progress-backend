package main

import (
	"log"
	"net/http"

	"github.com/thewizardplusplus/go-upload-progress/gateways/handlers"
)

func main() {
	http.HandleFunc("/api/v1/files", func(w http.ResponseWriter, r *http.Request) {
		var fileHandler handlers.FileHandler
		switch r.Method {
		case http.MethodGet:
			fileHandler.GetFiles(w, r)
		case http.MethodPost:
			fileHandler.SaveFile(w, r)
		case http.MethodDelete:
			if filename := r.FormValue("filename"); filename != "" {
				DeleteFileHandler(w, r)
			} else {
				DeleteFilesHandler(w, r)
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
