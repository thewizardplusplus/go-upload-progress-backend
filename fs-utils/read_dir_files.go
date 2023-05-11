package fsutils

import (
	"io/fs"
)

func ReadDirFiles(fsInstance fs.FS, path string) ([]fs.DirEntry, error) {
	dirEntries, err := fs.ReadDir(fsInstance, path)
	if err != nil {
		return nil, err
	}

	files := make([]fs.DirEntry, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			files = append(files, dirEntry)
		}
	}

	return files, nil
}
