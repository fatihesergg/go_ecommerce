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

type AuthHandler struct {
	UserService service.UserService
	Validator   *validator.Validate
}

func NewAuthHandler(service service.UserService, validator *validator.Validate) AuthHandler {
	return AuthHandler{UserService: service, Validator: validator}
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.Login	true	"credentials"
//	@Success		200			{object}	util.ApiResponse{}
//	@Failure		400			{object}	util.ApiResponse{}
//	@Failure		500			{object}	util.ApiResponse{}
//	@Router			/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.Login
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
	user, err := h.UserService.GetByEmail(data.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "User not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting user"
		util.WriteJson(w, response)
		return
	}
	var role string
	if user.Role == "admin" {
		role = "admin"
	} else {
		role = "user"
	}

	userIdint := strconv.Itoa(int(user.ID))
	token, err := util.CreateJWT(userIdint, role)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating jwt token"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = token
	util.WriteJson(w, response)
}

// Register godoc
//
//	@Summary		Register
//	@Description	Register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			register	body		dto.Register	true	"User informations"
//	@Success		200			{object}	util.ApiResponse{}
//	@Failure		400			{object}	util.ApiResponse{}
//	@Failure		500			{object}	util.ApiResponse{}
//	@Router			/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.Register
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
	exist, _ := h.UserService.GetByEmail(data.Email)
	if exist.Email == data.Email {
		response.Status = http.StatusBadRequest
		response.Message = "User already exist"
		util.WriteJson(w, response)
		return
	}
	user := model.User{
		Name:     data.Name,
		LastName: data.LastName,
		UserName: data.UserName,
		Email:    data.Email,
		Role:     "user",
		Password: data.Password,
	}
	err = h.UserService.Create(user)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating user."
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}
