package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/presentation/request"
	"app/api/presentation/response"
	"net/http"

	"github.com/pkg/errors"
)

type ThreadHandler interface {
	Create(w http.ResponseWriter, r *http.Request)               //Create thread
	GetAll(w http.ResponseWriter, r *http.Request)               //Get all threads
	GetByID(w http.ResponseWriter, r *http.Request)              //Get thread by ID
	GetOnlyPublic(w http.ResponseWriter, r *http.Request)        //Get public thread
	GetMembersByThreadID(w http.ResponseWriter, r *http.Request) //Get members in thread
	Update(w http.ResponseWriter, r *http.Request)               //Thread update
	Delete(w http.ResponseWriter, r *http.Request)               //Thread delete
	Join(w http.ResponseWriter, r *http.Request)                 //Join member to thread
	Leave(w http.ResponseWriter, r *http.Request)                //Leave the thread
	ForceToLeave(w http.ResponseWriter, r *http.Request)         //Kicked the member from thread
}

type threadHandler struct {
	threadInteractor interactor.ThreadInteractor
}

func NewThreadHandler(ti interactor.ThreadInteractor) ThreadHandler {
	return &threadHandler{
		threadInteractor: ti,
	}
}

func (th *threadHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	src, err := ReadRequestBody(r, &request.CreateThreadRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.CreateThreadRequest)
	err = req.Validation()
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	thread, err := th.threadInteractor.Create(req.Name, req.Description, req.LimitUsers, req.IsPublic, userID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create thread"), "failed to create thread")
		return
	}

	response.Success(w, response.ConvertToThreadResponse(thread))
}

func (th *threadHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	threads, err := th.threadInteractor.GetAll()
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get threads"), "failed to get threads")
		return
	}
	response.Success(w, response.ConvertToThreadsResponse(threads))
}

func (th *threadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	thread, err := th.threadInteractor.GetByID(id)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get thread"), "failed to get thread")
		return
	}
	response.Success(w, response.ConvertToThreadResponse(thread))
}

func (th *threadHandler) GetOnlyPublic(w http.ResponseWriter, r *http.Request) {
	threads, err := th.threadInteractor.GetOnlyPublic()
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get threads"), "failed to get threads")
		return
	}
	response.Success(w, response.ConvertToThreadsResponse(threads))
}

func (th *threadHandler) GetMembersByThreadID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}
	members, err := th.threadInteractor.GetMembersByThreadID(id)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get members"), "failed to get members")
		return
	}
	response.Success(w, response.ConvertToUsersResponse(members))
}

func (th *threadHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}
	id, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	src, err := ReadRequestBody(r, &request.UpdateThreadRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.UpdateThreadRequest)
	err = req.Validation(th.threadInteractor, id, userID)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	thread, err := th.threadInteractor.Update(id, req.Name, req.Description, req.LimitUsers, req.IsPublic)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to update thread"), "failed to update thread")
		return
	}
	response.Success(w, response.ConvertToThreadResponse(thread))
}

func (th *threadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := ReadPathParam(r, "id")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	if err := th.threadInteractor.Delete(id); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to delete thread"), "failed to delete thread")
		return
	}
	response.NoContent(w)
}

func (th *threadHandler) Join(w http.ResponseWriter, r *http.Request) {
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
	_, err = th.threadInteractor.GetByID(threadID)
	if err != nil {
		response.NotFound(w, errors.Wrap(err, "failed to get thread"), "thread is not found")
		return
	}

	if err = th.threadInteractor.AddMember(threadID, userID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to join thread"), "failed to join thread")
		return
	}
	response.NoContent(w)
}

func (th *threadHandler) Leave(w http.ResponseWriter, r *http.Request) {
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

	if err = th.threadInteractor.RemoveMember(threadID, userID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to leave thread"), "failed to leave thread")
		return
	}
	response.NoContent(w)
}

func (th *threadHandler) ForceToLeave(w http.ResponseWriter, r *http.Request) {
	response.NotImplemented(w)
}
