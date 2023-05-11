package entities

import (
	"io/fs"
)

type FileInfoGroup []FileInfo

func NewFileInfoGroup(files []fs.DirEntry) (FileInfoGroup, error) {
	fileInfos := make(FileInfoGroup, 0, len(files))
	for _, file := range files {
		fileInfo, err := NewFileInfoFromDirEntry(file)
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}
