package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/thewizardplusplus/go-upload-progress/entities"
)

type FileUsecase interface {
	GetFiles() ([]entities.FileInfo, error)
	SaveFile(file io.Reader, filename string) error
	DeleteFile(filename string) error
	DeleteFiles() error
}

type Logger interface {
	Print(values ...any)
	Printf(format string, args ...any)
}

type FileHandler struct {
	FileUsecase FileUsecase
	Logger      Logger
}

func (h FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetFiles(w, r)
	case http.MethodPost:
		h.SaveFile(w, r)
	case http.MethodDelete:
		if filename := r.FormValue("filename"); filename != "" {
			h.DeleteFile(w, r)
		} else {
			h.DeleteFiles(w, r)
		}
	default:
		errStatus := http.StatusMethodNotAllowed
		h.handleError(w, errors.New(http.StatusText(errStatus)), errStatus)
	}
}

func (h FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := h.FileUsecase.GetFiles()
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(fileInfos)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes) // nolint: errcheck
}

func (h FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.FileUsecase.SaveFile(file, fileHeader.Filename); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	if filename == "" {
		h.handleError(w, errors.New("filename is required"), http.StatusBadRequest)
		return
	}

	if err := h.FileUsecase.DeleteFile(filename); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h FileHandler) DeleteFiles(w http.ResponseWriter, r *http.Request) {
	if err := h.FileUsecase.DeleteFiles(); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h FileHandler) handleError(w http.ResponseWriter, err error, status int) {
	h.Logger.Print(err)
	http.Error(w, err.Error(), status)
}
