package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func SaveFileHandler(w http.ResponseWriter, r *http.Request) {
	sourceFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	generatedFilename, err := GenerateFilename("./files", fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetFile, err := os.Create(generatedFilename)
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

func GenerateFilename(baseDir string, originalFilename string) (string, error) {
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

	generatedFilename := originalFilename
	for {
		if _, exists := filenameSet[generatedFilename]; !exists {
			break
		}

		extension := filepath.Ext(generatedFilename)
		randomSuffix := strconv.Itoa(rand.Int())
		generatedFilename = strings.TrimSuffix(generatedFilename, extension) +
			randomSuffix +
			extension
	}

	fullGeneratedFilename := filepath.Join(baseDir, generatedFilename)
	return fullGeneratedFilename, nil
}
