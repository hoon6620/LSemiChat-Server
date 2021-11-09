package handler

import (
	"app/api/application/interactor"
	"app/api/domain/entity"
	"app/api/infrastructure/lcontext"
	"app/api/presentation/request"
	"app/api/presentation/response"
	"net/http"

	"github.com/pkg/errors"
)

type MessageHandler interface {
	Create(w http.ResponseWriter, r *http.Request)        //Create Massage
	GetByThreadID(w http.ResponseWriter, r *http.Request) //Get Thread Messages
	AddFavorite(w http.ResponseWriter, r *http.Request)   //Add favorite(Like) message
}

type messageHandler struct {
	messageInteractor interactor.MessageInteractor
	threadInteractor  interactor.ThreadInteractor
}

func NewMessageHandler(mi interactor.MessageInteractor, ti interactor.ThreadInteractor) MessageHandler {
	return &messageHandler{
		messageInteractor: mi,
		threadInteractor:  ti,
	}
}

func (mh *messageHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	src, err := ReadRequestBody(r, &request.CreateMessageRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}

	req, _ := src.(*request.CreateMessageRequest)
	err = req.Validation()
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	message, err := mh.messageInteractor.Create(req.Message, req.Grade, userID, threadID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create message"), "failed to create message")
		return
	}
	response.Success(w, response.ConvertToMessageResponse(message))
}

func (mh *messageHandler) GetByThreadID(w http.ResponseWriter, r *http.Request) {
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

	members, err := mh.threadInteractor.GetMembersByThreadID(threadID)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to members of thread"), "failed to members of thread")
		return
	}

	if !checkMember(userID, members) {
		response.BadRequest(w, errors.New("not memeber of thread"), "not member of thread")
		return
	}

	messages, err := mh.messageInteractor.GetByThreadID(threadID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get messages"), "failed to get messages")
		return
	}
	response.Success(w, response.ConvertToMessagesResponse(messages))
}

func (mh *messageHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	threadID, err := ReadPathParam(r, "threadID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	messageID, err := ReadPathParam(r, "messageID")
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "path parameter is empty"), "path parameter is empty")
		return
	}

	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	_, err = mh.threadInteractor.GetByID(threadID)
	if err != nil {
		response.NotFound(w, errors.Wrap(err, "failed to find thread"), "failed to find thread")
		return
	}
	_, err = mh.messageInteractor.GetByID(messageID)
	if err != nil {
		response.NotFound(w, errors.Wrap(err, "failed to find message"), "failed to find message")
		return
	}

	if err = mh.messageInteractor.AddFavorite(messageID, userID); err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to add favorite"), "failed to add favorite")
		return
	}
	response.NoContent(w)
}

func checkMember(userID string, members []*entity.User) bool {
	for _, member := range members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}
