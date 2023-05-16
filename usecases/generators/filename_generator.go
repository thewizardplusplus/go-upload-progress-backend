package generators

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type FilenameGenerator struct {
	MaximumNumberOfTries  int
	RandomSuffixByteCount int
}

func (g FilenameGenerator) GenerateUniqueFilename(
	baseFilename string,
	existingFiles []fs.DirEntry,
) (string, error) {
	existingFilenameSet := make(map[string]struct{})
	for _, existingFile := range existingFiles {
		existingFilenameSet[existingFile.Name()] = struct{}{}
	}

	uniqueFilename := baseFilename
	for try := 0; ; try++ {
		if _, exists := existingFilenameSet[uniqueFilename]; !exists {
			break
		}

		if try >= g.MaximumNumberOfTries {
			return "", errors.New("the maximum number of tries is reached")
		}

		var err error
		uniqueFilename, err = g.GenerateRandomFilename(baseFilename)
		if err != nil {
			return "", err
		}
	}

	return uniqueFilename, nil
}

func (g FilenameGenerator) GenerateRandomFilename(
	baseFilename string,
) (string, error) {
	randomSuffixBytes := make([]byte, g.RandomSuffixByteCount)
	if _, err := rand.Read(randomSuffixBytes); err != nil {
		return "", err
	}

	randomSuffix := fmt.Sprintf("_%x", randomSuffixBytes)
	return addSuffixToFilename(baseFilename, randomSuffix), nil
}

func addSuffixToFilename(filename string, suffix string) string {
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension) + suffix + extension
}
