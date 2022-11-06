package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func DeleteFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullFilename := filepath.Join("./files", file.Name())
		if err := os.Remove(fullFilename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
