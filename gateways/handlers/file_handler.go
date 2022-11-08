package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/thewizardplusplus/go-upload-progress-backend/entities"
)

type FileUsecase interface {
	GetFiles() ([]entities.FileInfo, error)
	SaveFile(file io.Reader, filename string) error
}

type FileHandler struct {
	FileUsecase FileUsecase
}

func (h FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := h.FileUsecase.GetFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(fileInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes) // nolint: errcheck
}

func (h FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.FileUsecase.SaveFile(file, fileHeader.Filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
