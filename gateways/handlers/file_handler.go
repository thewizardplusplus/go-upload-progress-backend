package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thewizardplusplus/go-upload-progress/entities"
)

const (
	randomSuffixByteCount = 4
)

type FileHandler struct{}

func (h FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileInfos := make([]entities.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileInfos = append(fileInfos, entities.NewFileInfo(fileInfo))
	}

	sort.Slice(fileInfos, func(i int, j int) bool {
		return fileInfos[i].ModTime.After(fileInfos[j].ModTime) // reverse order
	})

	bytes, err := json.Marshal(fileInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes) // nolint: errcheck
}

func (h FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	sourceFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	uniqueFilename, err := generateUniqueFilename("./files", fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetFile, err := os.Create(uniqueFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func generateUniqueFilename(baseDir string, originalFilename string) (string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}

	filenameSet := make(map[string]struct{})
	for _, file := range files {
		if !file.IsDir() {
			filenameSet[file.Name()] = struct{}{}
		}
	}

	uniqueFilename := originalFilename
	for {
		if _, exists := filenameSet[uniqueFilename]; !exists {
			break
		}

		randomSuffixBytes := make([]byte, randomSuffixByteCount)
		if _, err := rand.Read(randomSuffixBytes); err != nil {
			return "", err
		}

		extension := filepath.Ext(uniqueFilename)
		uniqueFilename = strings.TrimSuffix(uniqueFilename, extension) +
			hex.EncodeToString(randomSuffixBytes) +
			extension
	}

	return filepath.Join(baseDir, uniqueFilename), nil
}
