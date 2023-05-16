package writablefs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	writablefs "github.com/thewizardplusplus/go-writable-fs"
)

type DirFS struct {
	writablefs.DummyFS

	innerDirFS fs.FS
	baseDir    string
}

func NewDirFS(baseDir string) DirFS {
	return DirFS{
		innerDirFS: os.DirFS(baseDir),
		baseDir:    baseDir,
	}
}

func (dfs DirFS) Open(path string) (fs.File, error) {
	return dfs.innerDirFS.Open(path)
}

func (dfs DirFS) Stat(path string) (fs.FileInfo, error) {
	return dfs.innerDirFS.(fs.StatFS).Stat(path)
}

// Method `CreateExcl()` acts by analogy with function `os.Create()`,
// but replaces flag `os.O_TRUNC` with `os.O_EXCL`.
func (dfs DirFS) CreateExcl(path string) (writablefs.WritableFile, error) {
	// use the "open" operation, since the `os.Create()` uses it
	if err := checkPath(path, "open"); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		dfs.joinWithBaseDir(path),
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

func (dfs DirFS) Remove(path string) error {
	if err := checkPath(path, "remove"); err != nil {
		return err
	}

	if err := os.Remove(dfs.joinWithBaseDir(path)); err != nil {
		// restore the original path instead of its joined version
		updatePathInPathError(err, path)

		return err
	}

	return nil
}

func (dfs DirFS) joinWithBaseDir(path string) string {
	return filepath.Join(dfs.baseDir, path)
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
