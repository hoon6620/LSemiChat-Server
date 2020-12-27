package handler

import (
	"app/api/application/interactor"
	"app/api/presentation/request"
	"app/api/presentation/response"
	"net/http"

	"github.com/pkg/errors"
)

type CategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	categoryInteractor interactor.CategoryInteractor
}

func NewCategoryHandler(ci interactor.CategoryInteractor) CategoryHandler {
	return &categoryHandler{
		categoryInteractor: ci,
	}
}

func (ch *categoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	src, err := ReadRequestBody(r, &request.CreateCategoryRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.CreateCategoryRequest)
	if err = req.Validate(); err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	category, err := ch.categoryInteractor.Create(req.Category)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create category"), "failed to create category")
		return
	}

	response.Success(w, response.ConvertToCategoryResponse(category))
}

func (ch *categoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := ReadPathParam(r, "id")
	category, err := ch.categoryInteractor.GetByID(id)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get category"), "failed to get category")
		return
	}
	response.Success(w, response.ConvertToCategoryResponse(category))
}

func (ch *categoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.categoryInteractor.GetAll()
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to get categories"), "failed to get categories")
		return
	}
	response.Success(w, response.ConvertToCategoriesResponse(categories))
}

func (ch *categoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := ReadPathParam(r, "id")
	src, err := ReadRequestBody(r, &request.UpdateCategoryRequest{})
	if err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to read request"), "failed to read request")
		return
	}
	req, _ := src.(*request.UpdateCategoryRequest)
	if err = req.Validate(); err != nil {
		response.BadRequest(w, errors.Wrap(err, "failed to validation"), err.Error())
		return
	}

	category, err := ch.categoryInteractor.Update(id, req.Category)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to create category"), "failed to create category")
		return
	}

	response.Success(w, response.ConvertToCategoryResponse(category))
}

func (ch *categoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := ReadPathParam(r, "id")
	err := ch.categoryInteractor.Delete(id)
	if err != nil {
		response.InternalServerError(w, errors.Wrap(err, "failed to delete category"), "failed to delete category")
		return
	}
	response.NoContent(w)
}
