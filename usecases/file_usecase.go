package usecases

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thewizardplusplus/go-upload-progress/entities"
)

const (
	randomSuffixByteCount = 4
)

type FileUsecase struct {
	FileDir string
	FS      fs.FS
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
		return fileInfos[i].ModTime.After(fileInfos[j].ModTime) // reverse order
	})

	return fileInfos, nil
}

func (u FileUsecase) SaveFile(file io.Reader, filename string) error {
	uniqueFilename, err := u.makeUniqueFilename(filename)
	if err != nil {
		return err
	}

	savedFile, err := os.Create(filepath.Join(u.FileDir, uniqueFilename))
	if err != nil {
		return err
	}
	defer savedFile.Close()

	if _, err := io.Copy(savedFile, file); err != nil {
		return err
	}

	return nil
}

func (u FileUsecase) DeleteFile(filename string) error {
	return os.Remove(filepath.Join(u.FileDir, filename))
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
	dirEntries, err := fs.ReadDir(u.FS, ".")
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

		uniqueFilename, err = makeRandomFilename(uniqueFilename)
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
		hex.EncodeToString(randomSuffixBytes) +
		extension
	return randomFilename, nil
}
