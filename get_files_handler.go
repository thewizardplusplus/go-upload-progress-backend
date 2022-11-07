package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"

	"github.com/thewizardplusplus/go-upload-progress/entities"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
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
