package entities

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	Name    string
	SizeInB int64
	ModTime time.Time `format:"date-time"`
}

func NewFileInfo(fileInfo fs.FileInfo) FileInfo {
	return FileInfo{
		Name:    fileInfo.Name(),
		SizeInB: fileInfo.Size(),
		ModTime: fileInfo.ModTime(),
	}
}
