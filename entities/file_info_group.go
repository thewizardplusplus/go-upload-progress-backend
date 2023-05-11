package entities

import (
	"io/fs"
)

type FileInfoGroup []FileInfo

func NewFileInfoGroup(files []fs.DirEntry) (FileInfoGroup, error) {
	fileInfos := make(FileInfoGroup, 0, len(files))
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, NewFileInfo(fileInfo))
	}

	return fileInfos, nil
}
