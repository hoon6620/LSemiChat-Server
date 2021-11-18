package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/presentation/response"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type FileHandler interface {
	Download(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
}

type fileHandler struct {
	userInteractor   interactor.UserInteractor
	threadInteractor interactor.ThreadInteractor
	imgPath          string
}

func NewFileHandler(ui interactor.UserInteractor, ti interactor.ThreadInteractor) FileHandler {
	return &fileHandler{
		userInteractor:   ui,
		threadInteractor: ti,
		imgPath:          os.Getenv("FILE_PATH"),
	}
}

func (fh *fileHandler) Download(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "threadID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	if !fh.threadInteractor.IsParticipated(threadID, userID) {
		response.BadRequest(w, errors.New(userID+" are not participated in room "+threadID), userID+" are not participated in room "+threadID)
	}

	fileID, err := ReadPathParam(r, "fileID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	http.ServeFile(w, r, fh.imgPath+threadID+"/"+fileID)
}

func (fh *fileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "threadID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	if !fh.threadInteractor.IsParticipated(threadID, userID) {
		response.BadRequest(w, errors.New(userID+" are not participated in room "+threadID), userID+" are not participated in room "+threadID)
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(fh.imgPath + handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.Write(fileBytes)
}
