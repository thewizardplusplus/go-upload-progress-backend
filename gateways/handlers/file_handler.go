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
	SaveFile(file io.Reader, filename string) (entities.FileInfo, error)
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

// # swag tool annotations
//
//	@router /files/ [GET]
//	@summary Get files
//	@produce json
//	@success 200 {array} entities.FileInfo
//	@failure 500 {string} string
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

// # swag tool annotations
//
//	@router /files/ [POST]
//	@summary Save a file
//	@param file formData file true "file"
//	@accept multipart/form-data
//	@produce json
//	@success 201 {object} entities.FileInfo
//	@failure 400 {string} string
//	@failure 500 {string} string
func (h FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	savedFileInfo, err := h.FileUsecase.SaveFile(file, fileHeader.Filename)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(savedFileInfo)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes) // nolint: errcheck
}

// # swag tool annotations
//
// The described route is implemented via this method
// and the `DeleteFiles()` method.
//
//	@router /files/ [DELETE]
//	@summary Delete a file or files
//	@description If the filename is passed, the route will remove one file,
//	@description otherwise all files.
//	@param filename query string false "filename" minLength(1)
//	@produce plain
//	@success 204 {string} string
//	@failure 500 {string} string
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

// # swag tool annotations
//
// This method is used in the same route
// as the `DeleteFile()` method; see its comment.
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
