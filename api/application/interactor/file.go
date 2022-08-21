package interactor

import "app/api/domain/service"

type FileInteractor interface {
	SaveFile(threadID string, fileName string, file []byte) (string, error)
	LoadFile(threadID string, fileName string) ([]byte, error)
	SetUserIcon(userID string, file []byte) error
	GetUserIcon(userID string) ([]byte, error)
	SetThreadIcon(threadID string, file []byte) error
	GetThreadIcon(threadID string) ([]byte, error)
}

type fileInteractor struct {
	fileService service.FileService
}

func NewFileInteractor(fs service.FileService) FileInteractor {
	return &fileInteractor{
		fileService: fs,
	}
}

func (fi *fileInteractor) SaveFile(threadID string, fileName string, file []byte) (string, error) {
	return fi.fileService.SaveFile(threadID, fileName, file)
}

func (fi *fileInteractor) LoadFile(threadID string, fileName string) ([]byte, error) {
	return fi.fileService.LoadFile(threadID, fileName)
}

func (fi *fileInteractor) SetUserIcon(userID string, file []byte) error {
	return fi.fileService.SetUserIcon(userID, file)
}

func (fi *fileInteractor) GetUserIcon(userID string) ([]byte, error) {
	return fi.fileService.GetUserIcon(userID)
}

func (fi *fileInteractor) SetThreadIcon(threadID string, file []byte) error {
	return fi.fileService.SetThreadIcon(threadID, file)
}

func (fi *fileInteractor) GetThreadIcon(threadID string) ([]byte, error) {
	return fi.fileService.GetThreadIcon(threadID)
}
