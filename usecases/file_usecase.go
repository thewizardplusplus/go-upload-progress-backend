package usecases

import (
	"encoding/hex"
	"io"
	"io/fs"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thewizardplusplus/go-upload-progress/entities"
	"github.com/thewizardplusplus/go-upload-progress/gateways/writablefs"
)

const (
	randomSuffixByteCount = 4
)

type WritableFS interface {
	fs.FS

	Create(filename string) (writablefs.WritableFile, error)
	Remove(filename string) error
}

type FileUsecase struct {
	WritableFS WritableFS
}

func (u FileUsecase) GetFiles() ([]entities.FileInfo, error) {
	files, err := u.readDirFiles()
	if err != nil {
		return nil, err
	}

	fileInfos := make([]entities.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, entities.NewFileInfo(fileInfo))
	}

	sort.Slice(fileInfos, func(i int, j int) bool {
		return fileInfos[i].ModificationTime.
			After(fileInfos[j].ModificationTime) // reverse order
	})

	return fileInfos, nil
}

func (u FileUsecase) SaveFile(file io.Reader, filename string) (entities.FileInfo, error) {
	uniqueFilename, err := u.makeUniqueFilename(filename)
	if err != nil {
		return entities.FileInfo{}, err
	}

	savedFile, err := u.WritableFS.Create(uniqueFilename)
	if err != nil {
		return entities.FileInfo{}, err
	}
	defer savedFile.Close()

	if _, err := io.Copy(savedFile, file); err != nil {
		return entities.FileInfo{}, err
	}

	savedFileInfo, err := savedFile.Stat()
	if err != nil {
		return entities.FileInfo{}, err
	}

	return entities.NewFileInfo(savedFileInfo), nil
}

func (u FileUsecase) DeleteFile(filename string) error {
	return u.WritableFS.Remove(filename)
}

func (u FileUsecase) DeleteFiles() error {
	files, err := u.readDirFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := u.DeleteFile(file.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (u FileUsecase) readDirFiles() ([]fs.DirEntry, error) {
	dirEntries, err := fs.ReadDir(u.WritableFS, ".")
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

func (u FileUsecase) makeUniqueFilename(filename string) (string, error) {
	files, err := u.readDirFiles()
	if err != nil {
		return "", err
	}

	filenameSet := make(map[string]struct{})
	for _, file := range files {
		filenameSet[file.Name()] = struct{}{}
	}

	uniqueFilename := filename
	for {
		if _, exists := filenameSet[uniqueFilename]; !exists {
			break
		}

		uniqueFilename, err = makeRandomFilename(filename)
		if err != nil {
			return "", err
		}
	}

	return uniqueFilename, nil
}

func makeRandomFilename(filename string) (string, error) {
	randomSuffixBytes := make([]byte, randomSuffixByteCount)
	if _, err := rand.Read(randomSuffixBytes); err != nil {
		return "", err
	}

	extension := filepath.Ext(filename)
	randomFilename := strings.TrimSuffix(filename, extension) +
		"_" + hex.EncodeToString(randomSuffixBytes) +
		extension
	return randomFilename, nil
}
