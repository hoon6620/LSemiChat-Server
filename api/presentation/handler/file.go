package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/llog"
	"app/api/presentation/response"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type FileHandler interface {
	Download(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	GetUserIcon(w http.ResponseWriter, r *http.Request)
	SetUserIcon(w http.ResponseWriter, r *http.Request)
	GetThreadIcon(w http.ResponseWriter, r *http.Request)
	SetThreadIcon(w http.ResponseWriter, r *http.Request)
}

type fileHandler struct {
	fileInteractor    interactor.FileInteractor
	userInteractor    interactor.UserInteractor
	threadInteractor  interactor.ThreadInteractor
	messageInteractor interactor.MessageInteractor
}

func NewFileHandler(fi interactor.FileInteractor, ui interactor.UserInteractor, ti interactor.ThreadInteractor, mi interactor.MessageInteractor) FileHandler {
	return &fileHandler{
		fileInteractor:    fi,
		userInteractor:    ui,
		threadInteractor:  ti,
		messageInteractor: mi,
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

	file, err := fh.fileInteractor.LoadFile(threadID, fileID)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "filed to get file"), "filed to get file")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file)
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
		response.Unauthorized(w, errors.Wrap(err, "'uploadFile' is not contined"), "uploadFile is not contined")
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to read file"), "failed to read file")
		return
	}

	fileName, err := fh.fileInteractor.SaveFile(threadID, handler.Filename, fileBytes)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to save file"), "failed to save file")
		return
	}

	message, err := fh.messageInteractor.Create(fileName, 10, userID, threadID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create message"), "failed to create message")
		return
	}

	response.Success(w, response.ConvertToMessageResponse(message))
}

func (fh *fileHandler) SetUserIcon(w http.ResponseWriter, r *http.Request) {
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
	fh.fileInteractor.SetUserIcon(userID, fileBytes)

	response.Success(w, &response.FileResponse{
		FilePath: "/users/" + userID + "/icon",
	})
}

func (fh *fileHandler) GetUserIcon(w http.ResponseWriter, r *http.Request) {
	userID, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	file, err := fh.fileInteractor.GetUserIcon(userID)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "filed to get user icon"), "filed to get user icon")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(file)
}

func (fh *fileHandler) SetThreadIcon(w http.ResponseWriter, r *http.Request) {
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}
	threadID, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	thread, err := fh.threadInteractor.GetByID(threadID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get thread"), "failed to get thread")
		return
	}
	user, err := fh.userInteractor.GetByUserID(userID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get user"), "failed to get user")
		return
	}

	if thread.Author.ID != user.ID {
		response.BadRequest(w, errors.New("do not have permissions to change icon"), "do not have permissions to change room icon")
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("threadIcon")
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

	err = fh.fileInteractor.SetThreadIcon(threadID, fileBytes)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to save file"), "failed to save file")
		return
	}

	response.Success(w, &response.FileResponse{
		FilePath: "/threads/" + threadID + "/icon",
	})
}

func (fh *fileHandler) GetThreadIcon(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	file, err := fh.fileInteractor.GetThreadIcon(threadID)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "filed to get user icon"), "filed to get user icon")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(file)
}
