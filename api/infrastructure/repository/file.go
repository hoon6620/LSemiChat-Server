package repository

import (
	"app/api/constants"
	"app/api/domain/repository"
	"app/api/llog"
	"io/ioutil"
	"os"
)

type fileRepository struct {
	filePath string
}

func NewFileRepository() repository.FileRepository {
	//create image dir
	imgPath := os.Getenv("FILE_PATH")
	if imgPath == "" {
		imgPath = constants.ImgPath
	}
	createDirectory(imgPath)
	createDirectory(imgPath + "/users")
	createDirectory(imgPath + "/threads")

	return &fileRepository{
		filePath: imgPath,
	}
}

func (fr *fileRepository) CreateFile(threadID string, fileName string, file []byte) error {
	path := fr.filePath + "/threads/" + threadID
	return writeFile(path+"/"+fileName, file)
}

func (fr *fileRepository) GetFile(threadID string, fileName string) ([]byte, error) {
	return readFile(fr.filePath + "/threads/" + threadID + "/" + fileName)
}

func (fr *fileRepository) SetUserIcon(userID string, file []byte) error {
	return writeFile(fr.filePath+"/users/"+userID+"/icon.jpg", file)
}

func (fr *fileRepository) GetUserIcon(userID string) ([]byte, error) {
	return readFile(fr.filePath + "/users/" + userID + "/icon.jpg")

}

func (fr *fileRepository) SetThreadIcon(threadID string, file []byte) error {
	return writeFile(fr.filePath+"/threads/"+threadID+"/icon.jpg", file)
}

func (fr *fileRepository) GetThreadIcon(threadID string) ([]byte, error) {
	return readFile(fr.filePath + "/threads/" + threadID + "/icon.jpg")
}

func (fr *fileRepository) CreateThreadDir(threadID string) error {
	return createDirectory(fr.filePath + "/threads/" + threadID)

}

func (fr *fileRepository) CreateUserDir(userID string) error {
	return createDirectory(fr.filePath + "/users/" + userID)
}

func createDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filePath string, file []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	llog.Debug(filePath)
	llog.Debug(len(file))

	f.Write(file)
	return nil
}

func readFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
