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

// Method `CreateExcl()` acts by analogy with function `os.Create()`,
// but replaces flag `os.O_TRUNC` with `os.O_EXCL`.
func (wfs WritableFS) CreateExcl(filename string) (WritableFile, error) {
	// use the "open" operation, since the `os.Create()` uses it
	if err := checkFilename(filename, "open"); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		wfs.joinWithBaseDir(filename),
		os.O_RDWR|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		// restore the original filename instead of its joined version
		updatePathInPathError(err, filename)

		return nil, err
	}

	return file, nil
}

func (wfs WritableFS) Remove(filename string) error {
	if err := checkFilename(filename, "remove"); err != nil {
		return err
	}

	if err := os.Remove(wfs.joinWithBaseDir(filename)); err != nil {
		// restore the original filename instead of its joined version
		updatePathInPathError(err, filename)

		return err
	}

	return nil
}

func (wfs WritableFS) joinWithBaseDir(filename string) string {
	return filepath.Join(wfs.baseDir, filename)
}

// This check is made for consistency with the implementation of `os.DirFS()`.
//
// # License
//
//	BSD 3-Clause "New" or "Revised" License
//	Copyright (C) 2009 The Go Authors
func checkFilename(filename string, operation string) error {
	if !fs.ValidPath(filename) ||
		(runtime.GOOS == "windows" && strings.ContainsAny(filename, `\:`)) {
		return &fs.PathError{Op: operation, Path: filename, Err: fs.ErrInvalid}
	}

	return nil
}

func updatePathInPathError(err error, path string) {
	err.(*fs.PathError).Path = path
}
