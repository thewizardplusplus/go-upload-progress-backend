package entities

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64     // in bytes
	ModTime time.Time `format:"date-time"`
}

func NewFileInfo(fileInfo fs.FileInfo) FileInfo {
	return FileInfo{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		ModTime: fileInfo.ModTime(),
	}
}
