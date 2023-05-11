package usecases

// nolint: lll
import (
	"io"
	"io/fs"
	"sort"

	"github.com/thewizardplusplus/go-upload-progress-backend/entities"
	writablefs "github.com/thewizardplusplus/go-upload-progress-backend/gateways/writable-fs"
)

type WritableFS interface {
	fs.FS

	Create(filename string) (writablefs.WritableFile, error)
	Remove(filename string) error
}

type FilenameGenerator interface {
	GenerateUniqueFilename(
		baseFilename string,
		existingFiles []fs.DirEntry,
	) (string, error)
}

type FileUsecase struct {
	WritableFS        WritableFS
	FilenameGenerator FilenameGenerator
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

func (u FileUsecase) SaveFile(
	file io.Reader,
	filename string,
) (entities.FileInfo, error) {
	existingFiles, err := u.readDirFiles()
	if err != nil {
		return entities.FileInfo{}, err
	}

	uniqueFilename, err :=
		u.FilenameGenerator.GenerateUniqueFilename(filename, existingFiles)
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

	return entities.NewFileInfoFromFile(savedFile)
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
