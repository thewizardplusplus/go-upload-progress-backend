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
