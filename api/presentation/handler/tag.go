package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/presentation/request"
	"app/api/presentation/response"
	"net/http"

	"github.com/pkg/errors"
)

type TagHandler interface {
	Create(w http.ResponseWriter, r *http.Request)              //Creeate tag
	GetAll(w http.ResponseWriter, r *http.Request)              //Get all tags
	GetByID(w http.ResponseWriter, r *http.Request)             //Get tag by ID
	AddTagToUser(w http.ResponseWriter, r *http.Request)        //Add tag to user
	RemoveTagFromUser(w http.ResponseWriter, r *http.Request)   //Remove tag from user
	AddTagToThread(w http.ResponseWriter, r *http.Request)      //Add tag to thread
	RemoveTagFromThread(w http.ResponseWriter, r *http.Request) //Remove tag from thread
}

type tagHandler struct {
	tagInteractor      interactor.TagInteractor
	categoryInteractor interactor.CategoryInteractor
}

func NewTagHandler(ti interactor.TagInteractor, ci interactor.CategoryInteractor) TagHandler {
	return &tagHandler{
		tagInteractor:      ti,
		categoryInteractor: ci,
	}
}

func (th *tagHandler) Create(w http.ResponseWriter, r *http.Request) {
	src, err := ReadRequestBody(r, &request.CreateTagRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.CreateTagRequest)
	err = req.Validation(th.categoryInteractor)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	tag, err := th.tagInteractor.Create(req.Tag, req.CategoryID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create tag"), "failed to create tag")
		return
	}

	response.Success(w, response.ConvertToTagResponse(tag))
}

func (th *tagHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tags, err := th.tagInteractor.GetAll()
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get"), "failed to get tags")
		return
	}
	response.Success(w, response.ConvertToTagsResponse(tags))
}

func (th *tagHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	tag, err := th.tagInteractor.GetByID(id)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get tag"), "failed to get tag")
		return
	}
	response.Success(w, response.ConvertToTagResponse(tag))
}

func (th *tagHandler) AddTagToUser(w http.ResponseWriter, r *http.Request) {
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	src, err := ReadRequestBody(r, &request.CreateTagRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.CreateTagRequest)
	err = req.Validation(th.categoryInteractor)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	if err = th.tagInteractor.AddToUser(req.Tag, req.CategoryID, userID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to add tag to user"), "failed to add tag to user")
		return
	}

	response.NoContent(w)
}

func (th *tagHandler) RemoveTagFromUser(w http.ResponseWriter, r *http.Request) {
	tagID, err := ReadPathParam(r, "tagID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}
	if err := th.tagInteractor.RemoveFromUser(tagID, userID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to remove tag from user"), "failed to remove tag from user")
		return
	}
	response.NoContent(w)
}

func (th *tagHandler) AddTagToThread(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "threadID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	src, err := ReadRequestBody(r, &request.CreateTagRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.CreateTagRequest)
	err = req.Validation(th.categoryInteractor)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	if err = th.tagInteractor.AddToThread(req.Tag, req.CategoryID, threadID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to add tag to thread"), "failed to add tag to thread")
		return
	}

	response.NoContent(w)
}

func (th *tagHandler) RemoveTagFromThread(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "threadID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	tagID, err := ReadPathParam(r, "tagID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	if err := th.tagInteractor.RemoveFromThread(tagID, threadID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to remove tag from user"), "failed to remove tag from user")
		return
	}
	response.NoContent(w)
}
