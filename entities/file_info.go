package entities

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}

func NewFileInfo(fileInfo fs.FileInfo) FileInfo {
	return FileInfo{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		ModTime: fileInfo.ModTime(),
	}
}
