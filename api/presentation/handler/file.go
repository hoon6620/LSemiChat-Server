package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/llog"
	"app/api/presentation/response"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type FileHandler interface {
	Download(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	GetUserIcon(w http.ResponseWriter, r *http.Request)
	SetUserIcon(w http.ResponseWriter, r *http.Request)
}

type fileHandler struct {
	threadInteractor  interactor.ThreadInteractor
	messageInteractor interactor.MessageInteractor
	imgPath           string
}

func NewFileHandler(ti interactor.ThreadInteractor, mi interactor.MessageInteractor) FileHandler {
	return &fileHandler{
		threadInteractor:  ti,
		messageInteractor: mi,
		imgPath:           os.Getenv("FILE_PATH"),
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
		llog.Info("test")
		//response.BadRequest(w, errors.New(userID+" are not participated in room "+threadID), userID+" are not participated in room "+threadID)
		//return
	}

	fileID, err := ReadPathParam(r, "fileID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	if fInfo, err := os.Stat(fh.imgPath + "/threads/" + threadID + "/" + fileID); err != nil {
		response.BadRequest(w, errors.Wrap(err, "file is not exist"), "file is not exist")
		return
	} else {
		llog.Info(fInfo)
		http.ServeFile(w, r, fh.imgPath+"/threads/"+threadID+"/"+fileID)
	}
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

	// if !fh.threadInteractor.IsParticipated(threadID, userID) {
	// 	response.BadRequest(w, errors.New(userID+" are not participated in room "+threadID), userID+" are not participated in room "+threadID)
	// }

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "'uploadFile' is not contined"), "uploadFile is not contined")
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to read file"), "failed to read file")
		return
	}

	message, err := fh.messageInteractor.Create(handler.Filename, 10, userID, threadID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create message"), "failed to create message")
		return
	}

	if _, err := os.Stat(fh.imgPath + "/" + threadID); os.IsNotExist(err) {
		err := os.Mkdir(fh.imgPath+"/"+threadID, os.ModeDir)
		if err != nil {
			llog.Error(err)
		}
	}

	f, err := os.Create(fh.imgPath + "/threads/" + threadID + "/" + message.ID)
	if err != nil {
		llog.Error(err)
		return
	}
	defer f.Close()

	f.Write(fileBytes)

	response.Success(w, response.ConvertToMessageResponse(message))
}

func (fh *fileHandler) GetUserIcon(w http.ResponseWriter, r *http.Request) {
	userID, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	if fInfo, err := os.Stat(fh.imgPath + "/users/" + userID + "/icon"); err != nil {
		response.BadRequest(w, errors.Wrap(err, "file is not exist"), "file is not exist")
		return
	} else {
		llog.Info(fInfo)
		http.ServeFile(w, r, fh.imgPath+"/users/"+userID+"/icon")
	}
}

func (fh *fileHandler) SetUserIcon(w http.ResponseWriter, r *http.Request) {
	llog.Info("set user icon")
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("userIcon")
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to read form"), "failed to read form")
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to read file"), "failed to read file")
		return
	}
	f, err := os.Create(fh.imgPath + "/users/" + userID + "/icon")
	if err != nil {
		llog.Error(err)
		return
	}
	defer f.Close()

	llog.Info(f.Name())
	f.Write(fileBytes)
	response.Success(w, &response.FileResponse{
		FilePath: "/users/" + userID + "/img",
	})
}
