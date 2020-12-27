package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lsession"
	"app/api/presentation/request"
	"app/api/presentation/response"
	"net/http"

	"github.com/pkg/errors"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authInteractor interactor.AuthInteractor
}

func NewAuthHandler(ai interactor.AuthInteractor) AuthHandler {
	return &authHandler{
		authInteractor: ai,
	}
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	src, err := ReadRequestBody(r, &request.LoginRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.LoginRequest)

	// validation
	if req.UserID == "" || req.Password == "" {
		response.BadRequest(w, errors.New("request field is empty"), "failed to validation")
		return
	}

	err = ah.authInteractor.Login(req.UserID, req.Password)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to authentication"), "failed to authentication")
		return
	}

	token, err := lsession.StartSession(w, req.UserID)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to start session"), "failed to login")
		return
	}

	res := &response.LoginResponse{
		Token: token,
	}
	response.Success(w, res)
}

func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := lsession.EndSession(w, r)
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to logout"), "failed to logout")
		return
	}
	response.NoContent(w)
}
