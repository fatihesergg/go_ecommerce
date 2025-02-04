package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fatihesergg/go_ecommerce/internal/dto"
	"github.com/fatihesergg/go_ecommerce/internal/model"
	"github.com/fatihesergg/go_ecommerce/internal/service"
	"github.com/fatihesergg/go_ecommerce/internal/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
	Validator       *validator.Validate
}

func NewCategoryHandler(service service.CategoryService, validator *validator.Validate) CategoryHandler {
	return CategoryHandler{CategoryService: service, Validator: validator}
}

// Get godoc
//
//	@Tags			category
//	@Summary		Show a category
//	@Description	get category by ID
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	util.ApiResponse{data=model.Category}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/category/{id} [get]
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	var response util.ApiResponse
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid category id"
		util.WriteJson(w, response)
		return
	}
	category, err := h.CategoryService.Get(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Category not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting category"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = category
	util.WriteJson(w, response)
}

// GetAll godoc
//
//	@Tags			category
//	@Summary		Show all category
//
//	@Description	get all category
//	@Produce		json
//	@Success		200	{object}	util.ApiResponse{data=[]model.Category}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/category [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryService.GetAll()
	var response util.ApiResponse

	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting categories."
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = categories
	util.WriteJson(w, response)
}

// Create godoc
//
//	@Tags			category
//	@Summary		Create a category
//
//	@Description	Create  a category
//
//	@Accept			json
//
//	@Produce		json
//
//	@Security		Bearer
//	@Param			category	body		dto.CategoryCreateDto	true	"Create Category"
//	@Success		200			{object}	util.ApiResponse{}
//	@Failure		400			{object}	util.ApiResponse{}
//	@Failure		500			{object}	util.ApiResponse{}
//	@Router			/category [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data dto.CategoryCreateDto
	var response util.ApiResponse
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = util.JsonDecodeError.Error()
		util.WriteJson(w, response)
		return
	}

	err := h.Validator.Struct(data)
	if err != nil {
		ve := err.(validator.ValidationErrors)
		response.Status = http.StatusBadRequest
		response.Message = util.GetErrorMessages(ve)
		util.WriteJson(w, response)
		return
	}

	err = h.CategoryService.Create(data.Name)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating category"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusCreated
	response.Message = "Category created successfully."
	util.WriteJson(w, response)
}

// Update godoc
//
//	@Tags			category
//	@Summary		Update a category
//
//	@Description	Update  a category
//
//	@Accept			json
//
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			category	body		dto.CategoryUpdateDto	true	"Update Category"
//
//	@Success		200			{object}	util.ApiResponse{}
//	@Failure		400			{object}	util.ApiResponse{}
//	@Failure		500			{object}	util.ApiResponse{}
//	@Router			/category [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.CategoryUpdateDto
	var response util.ApiResponse
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = util.JsonDecodeError.Error()
		util.WriteJson(w, response)
		return
	}
	err := h.Validator.Struct(data)
	if err != nil {
		ve := err.(validator.ValidationErrors)
		response.Status = http.StatusBadRequest
		response.Message = util.GetErrorMessages(ve)
		util.WriteJson(w, response)
		return
	}
	err = h.CategoryService.Update(model.Category{Name: data.Name, ID: uint(data.ID)})
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Category not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while updating category"
		util.WriteJson(w, response)
		return
	}

	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}
