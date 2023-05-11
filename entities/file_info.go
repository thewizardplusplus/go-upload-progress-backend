package entities

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	Name             string
	SizeInB          int64
	ModificationTime time.Time `format:"date-time"`
}

func NewFileInfo(fileInfo fs.FileInfo) FileInfo {
	return FileInfo{
		Name:             fileInfo.Name(),
		SizeInB:          fileInfo.Size(),
		ModificationTime: fileInfo.ModTime(),
	}
}

func NewFileInfoFromDirEntry(dirEntry fs.DirEntry) (FileInfo, error) {
	fileInfo, err := dirEntry.Info()
	if err != nil {
		return FileInfo{}, err
	}

	return NewFileInfo(fileInfo), nil
}

func NewFileInfoFromFile(file fs.File) (FileInfo, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return FileInfo{}, err
	}

	return NewFileInfo(fileInfo), nil
}
