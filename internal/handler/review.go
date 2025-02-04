package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fatihesergg/go_ecommerce/internal/dto"
	"github.com/fatihesergg/go_ecommerce/internal/middleware"
	"github.com/fatihesergg/go_ecommerce/internal/model"
	"github.com/fatihesergg/go_ecommerce/internal/service"
	"github.com/fatihesergg/go_ecommerce/internal/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	ReviewService  service.ReviewService
	UserService    service.UserService
	ProductService service.ProductService
	Validator      *validator.Validate
}

func NewReviewHandler(review service.ReviewService, user service.UserService, product service.ProductService, validator *validator.Validate) ReviewHandler {
	return ReviewHandler{ReviewService: review, UserService: user, ProductService: product, Validator: validator}
}

// Get godoc
//
//	@Tags			review
//	@Summary		Show a review
//	@Description	get review by ID
//	@Produce		json
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	util.ApiResponse{data=model.Review}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/review/{id} [get]
func (h *ReviewHandler) Get(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	var response util.ApiResponse
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid review id"
		util.WriteJson(w, response)
		return
	}

	review, err := h.ReviewService.Get(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Review not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting review"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = review
	util.WriteJson(w, response)
}

// Create godoc
//
//	@Tags			review
//	@Summary		Create a review
//	@Description	Create a review
//	@Produce		json
//	@Accept			json
//	@Security		BearerAuth
//	@Param			review	body		dto.ReviewCreateDto	true	"Review"
//	@Success		200		{object}	util.ApiResponse{}
//	@Failure		400		{object}	util.ApiResponse{}
//	@Failure		500		{object}	util.ApiResponse{}
//	@Router			/review [post]
func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data dto.ReviewCreateDto
	var response util.ApiResponse
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = util.JsonDecodeError.Error()
		util.WriteJson(w, response)
		return
	}

	if err := h.Validator.Struct(data); err != nil {
		ve := err.(validator.ValidationErrors)
		response.Status = http.StatusBadRequest
		response.Message = util.GetErrorMessages(ve)
		util.WriteJson(w, response)
		return
	}

	userID := r.Context().Value(middleware.AuthUserID).(string)
	user, err := h.UserService.Get(userID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	product, err := h.ProductService.Get(strconv.Itoa(data.ProductID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Product not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting product"
		util.WriteJson(w, response)
		return

	}
	review := model.Review{
		ProductID: product.ID,
		Product:   product,
		UserID:    user.ID,
		User:      user,
		Comment:   data.Comment,
	}
	err = h.ReviewService.Create(review)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating review"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}

// Update godoc
//
//	@Tags			review
//	@Summary		Update a review
//	@Description	Update a review
//	@Produce		json
//
//	@Accept			json
//	@Security		BearerAuth
//	@Param			review	body		dto.ReviewUpdateDto	true	"Review"
//	@Success		200		{object}	util.ApiResponse{}
//	@Failure		400		{object}	util.ApiResponse{}
//	@Failure		500		{object}	util.ApiResponse{}
//	@Router			/review [put]
func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.ReviewUpdateDto
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
	reqUserID := r.Context().Value(middleware.AuthUserID).(string)

	exist, err := h.ReviewService.Get(strconv.Itoa(data.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Review not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while updating review"
		util.WriteJson(w, response)
		return
	}
	reqUserIDint, _ := strconv.Atoi(reqUserID)
	if exist.UserID != uint(reqUserIDint) {
		response.Status = http.StatusBadRequest
		response.Message = "You have no permission to perform this action"
		util.WriteJson(w, response)
		return
	}

	review := model.Review{
		ID:        uint(data.ID),
		Comment:   data.Comment,
		UserID:    exist.UserID,
		User:      exist.User,
		ProductID: uint(data.ProductID),
		Product:   exist.Product,
	}

	err = h.ReviewService.Update(review)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Review not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while updating review"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}

// Delete godoc
//
//	@Tags			review
//	@Summary		Delete a review
//	@Description	Delete a review
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	util.ApiResponse{}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/review/{id} [delete]
func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	var response util.ApiResponse
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid review id"
		util.WriteJson(w, response)
		return
	}
	exist, err := h.ReviewService.Get(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Review not found"
			util.WriteJson(w, response)
			return
		}

		response.Status = http.StatusInternalServerError
		response.Message = "Error while deleting review"
		util.WriteJson(w, response)
		return

	}

	userID := r.Context().Value(middleware.AuthUserID).(string)
	if userID != strconv.Itoa(int(exist.UserID)) {
		response.Status = http.StatusBadRequest
		response.Message = "You have no permission to perform this action"
		util.WriteJson(w, response)
		return
	}
	err = h.ReviewService.Delete(r.PathValue("id"))
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while deleting review"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}
