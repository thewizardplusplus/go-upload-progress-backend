package writablefs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type WritableFile interface {
	fs.File
	io.Writer
}

type WritableFS struct {
	fs.FS

	baseDir string
}

func NewWritableFS(baseDir string) WritableFS {
	return WritableFS{
		FS:      os.DirFS(baseDir),
		baseDir: baseDir,
	}
}

func (writableFS WritableFS) Create(filename string) (WritableFile, error) {
	// use the "open" operation, since the `os.Create()` uses it
	if err := checkFilename(filename, "open"); err != nil {
		return nil, err
	}

	file, err := os.Create(writableFS.joinWithBaseDir(filename))
	if err != nil {
		// restore the original filename instead of its joined version
		err.(*fs.PathError).Path = filename

		return nil, err
	}

	return file, nil
}

func (writableFS WritableFS) Remove(filename string) error {
	if err := checkFilename(filename, "remove"); err != nil {
		return err
	}

	if err := os.Remove(writableFS.joinWithBaseDir(filename)); err != nil {
		// restore the original filename instead of its joined version
		err.(*fs.PathError).Path = filename

		return err
	}

	return nil
}

func (writableFS WritableFS) joinWithBaseDir(filename string) string {
	return filepath.Join(writableFS.baseDir, filename)
}

// This check is made for consistency with the implementation of `os.DirFS()`.
func checkFilename(filename string, operation string) error {
	// BSD 3-Clause "New" or "Revised" License
	// Copyright (C) 2009 The Go Authors

	if !fs.ValidPath(filename) ||
		(runtime.GOOS == "windows" && strings.ContainsAny(filename, `\:`)) {
		return &fs.PathError{Op: operation, Path: filename, Err: fs.ErrInvalid}
	}

	return nil
}
