package usecases

// nolint: lll
import (
	"io"
	"io/fs"
	"sort"

	"github.com/thewizardplusplus/go-upload-progress-backend/entities"
	fsutils "github.com/thewizardplusplus/go-upload-progress-backend/fs-utils"
	writablefs "github.com/thewizardplusplus/go-upload-progress-backend/gateways/writable-fs"
)

type WritableFS interface {
	fs.FS

	CreateExcl(filename string) (writablefs.WritableFile, error)
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

func (u FileUsecase) GetFiles() (entities.FileInfoGroup, error) {
	files, err := u.readTopLevelFSFiles()
	if err != nil {
		return nil, err
	}

	fileInfos, err := entities.NewFileInfoGroup(files)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(entities.ByModificationTime(fileInfos)))

	return fileInfos, nil
}

func (u FileUsecase) SaveFile(
	file io.Reader,
	filename string,
) (entities.FileInfo, error) {
	existingFiles, err := u.readTopLevelFSFiles()
	if err != nil {
		return entities.FileInfo{}, err
	}

	uniqueFilename, err :=
		u.FilenameGenerator.GenerateUniqueFilename(filename, existingFiles)
	if err != nil {
		return entities.FileInfo{}, err
	}

	savedFile, err := u.WritableFS.CreateExcl(uniqueFilename)
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
	files, err := u.readTopLevelFSFiles()
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

func (u FileUsecase) readTopLevelFSFiles() ([]fs.DirEntry, error) {
	return fsutils.ReadDirFiles(u.WritableFS, ".")
}
