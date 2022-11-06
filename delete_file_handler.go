package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	fullFilename := filepath.Join("./files", filename)
	if err := os.Remove(fullFilename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
