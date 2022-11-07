package usecases

import (
	"io/fs"
	"sort"

	"github.com/thewizardplusplus/go-upload-progress/entities"
)

type FileUsecase struct {
	FS fs.FS
}

func (u FileUsecase) GetFiles() ([]entities.FileInfo, error) {
	files, err := fs.ReadDir(u.FS, ".")
	if err != nil {
		return nil, err
	}

	fileInfos := make([]entities.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, entities.NewFileInfo(fileInfo))
	}

	sort.Slice(fileInfos, func(i int, j int) bool {
		return fileInfos[i].ModTime.After(fileInfos[j].ModTime) // reverse order
	})

	return fileInfos, nil
}
