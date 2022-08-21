package service

import (
	"app/api/domain/repository"
	"strings"
	"time"
)

type FileService interface {
	SaveFile(threadID string, fileName string, file []byte) (string, error)
	LoadFile(threadID string, fileName string) ([]byte, error)
	SetUserIcon(threadID string, file []byte) error
	GetUserIcon(threadID string) ([]byte, error)
	SetThreadIcon(threadID string, file []byte) error
	GetThreadIcon(threadID string) ([]byte, error)
}

type fileService struct {
	fileRepository repository.FileRepository
}

func NewFileService(fr repository.FileRepository) FileService {
	return &fileService{
		fileRepository: fr,
	}
}

func (fs *fileService) SaveFile(threadID string, fileName string, file []byte) (string, error) {
	name := createFileName(fileName)
	err := fs.fileRepository.CreateFile(threadID, name, file)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (fs *fileService) LoadFile(threadID string, fileName string) ([]byte, error) {
	return fs.fileRepository.GetFile(threadID, fileName)
}

func (fs *fileService) SetUserIcon(userID string, file []byte) error {
	return fs.fileRepository.SetUserIcon(userID, file)
}

func (fs *fileService) GetUserIcon(userID string) ([]byte, error) {
	return fs.fileRepository.GetUserIcon(userID)
}

func (fs *fileService) SetThreadIcon(threadID string, file []byte) error {
	return fs.fileRepository.SetThreadIcon(threadID, file)
}

func (fs *fileService) GetThreadIcon(threadID string) ([]byte, error) {
	return fs.fileRepository.GetThreadIcon(threadID)
}

func createFileName(fileName string) string {
	s := strings.Split(fileName, ".")
	return s[0] + time.Now().Format("20060102150405") + s[1]
}
