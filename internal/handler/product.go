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

type ProductHandler struct {
	CategoryService service.CategoryService
	ProductService  service.ProductService
	Validator       *validator.Validate
}

func NewProductHandler(productService service.ProductService, categoryService service.CategoryService, validator *validator.Validate) ProductHandler {
	return ProductHandler{ProductService: productService, CategoryService: categoryService, Validator: validator}
}

// Get godoc
//
//	@Tags			product
//	@Summary		Show a product
//	@Description	get product by ID
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	util.ApiResponse{data=model.Product}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/product/{id} [get]
func (h ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var response util.ApiResponse
	if id == "" {
		response.Status = http.StatusBadRequest
		response.Message = "Error getting id"
		util.WriteJson(w, response)
		return
	}
	product, err := h.ProductService.Get(id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Product not found"
			util.WriteJson(w, response)
			return
		}

		response.Status = http.StatusBadRequest
		response.Message = "Error getting product."
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = product
	util.WriteJson(w, response)
}

// GetAll godoc
//
//	@Tags			product
//	@Summary		Show all product
//	@Description	get products
//	@Produce		json
//	@Success		200	{object}	util.ApiResponse{data=[]model.Product}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/product [get]
func (h ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.GetAll()
	var response util.ApiResponse
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Product not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while gettig products"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = products
	util.WriteJson(w, response)
}

// Create godoc
//
//	@Tags			product
//	@Summary		Create a product
//
//	@Description	Create  a product
//
//	@Accept			json
//
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product	body		dto.ProductCreateDto	true	"Create Product"
//	@Success		200		{object}	util.ApiResponse{}
//	@Failure		400		{object}	util.ApiResponse{}
//	@Failure		500		{object}	util.ApiResponse{}
//	@Router			/product [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data dto.ProductCreateDto
	var response util.ApiResponse

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

	_, err = h.CategoryService.Get(strconv.Itoa(int(data.CategoryID)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Invalid category id"
			util.WriteJson(w, response)
			return
		}
	}
	product := model.Product{
		Name:       data.Name,
		ImageURL:   data.ImageURL,
		Price:      data.Price,
		CategoryID: data.CategoryID,
		Stock:      data.Stock,
	}
	err = h.ProductService.Create(product)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating product"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusCreated
	response.Message = "Product created succesfully."
	util.WriteJson(w, response)
}

// Update godoc
//
//	@Tags			product
//	@Summary		Update a product
//
//	@Description	Update a product
//
//	@Accept			json
//
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			product	body		dto.ProductUpdateDto	true	"Update Product"
//	@Success		200		{object}	util.ApiResponse{}
//	@Failure		400		{object}	util.ApiResponse{}
//	@Failure		500		{object}	util.ApiResponse{}
//	@Router			/product [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.ProductUpdateDto
	var response util.ApiResponse
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

	product := model.Product{ID: uint(data.ID), Name: data.Name, ImageURL: data.ImageURL, Price: data.Price, Stock: data.Stock, CategoryID: data.CategoryID}
	err = h.ProductService.Update(product)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Product not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while updating product."
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}

// Delete godoc
//
//	@Tags			product
//	@Summary		Delete a product
//
//	@Description	Delete  a product
//
//	@Accept			json
//
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			productID	path		int	true	"Product ID"
//	@Success		200			{object}	util.ApiResponse{}
//	@Failure		400			{object}	util.ApiResponse{}
//	@Failure		500			{object}	util.ApiResponse{}
//	@Router			/product/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := strconv.Atoi(id)
	var response util.ApiResponse
	if id == "" || err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid product id"
	}
	err = h.ProductService.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Product not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while deleting product"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}
