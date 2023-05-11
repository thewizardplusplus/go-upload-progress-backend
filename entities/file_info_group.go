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

type ByModificationTime FileInfoGroup

func (g ByModificationTime) Len() int {
	return len(g)
}

func (g ByModificationTime) Swap(i int, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g ByModificationTime) Less(i int, j int) bool {
	return g[i].ModificationTime.Before(g[j].ModificationTime)
}
