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
func (wfs WritableFS) CreateExcl(path string) (WritableFile, error) {
	// use the "open" operation, since the `os.Create()` uses it
	if err := checkPath(path, "open"); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		wfs.joinWithBaseDir(path),
		os.O_RDWR|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		// restore the original path instead of its joined version
		updatePathInPathError(err, path)

		return nil, err
	}

	return file, nil
}

func (wfs WritableFS) Remove(path string) error {
	if err := checkPath(path, "remove"); err != nil {
		return err
	}

	if err := os.Remove(wfs.joinWithBaseDir(path)); err != nil {
		// restore the original path instead of its joined version
		updatePathInPathError(err, path)

		return err
	}

	return nil
}

func (wfs WritableFS) joinWithBaseDir(path string) string {
	return filepath.Join(wfs.baseDir, path)
}

// This check is made for consistency with the implementation of `os.DirFS()`.
//
// # License
//
//	BSD 3-Clause "New" or "Revised" License
//	Copyright (C) 2009 The Go Authors
func checkPath(path string, operation string) error {
	if !fs.ValidPath(path) ||
		(runtime.GOOS == "windows" && strings.ContainsAny(path, `\:`)) {
		return &fs.PathError{Op: operation, Path: path, Err: fs.ErrInvalid}
	}

	return nil
}

func updatePathInPathError(err error, path string) {
	err.(*fs.PathError).Path = path
}
