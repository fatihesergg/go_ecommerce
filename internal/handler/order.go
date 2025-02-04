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

type OrderHandler struct {
	OrderService    service.OrderService
	ProductService  service.ProductService
	CategoryService service.CategoryService
	UserService     service.UserService
	Validator       *validator.Validate
}

func NewOrderHandler(orderService service.OrderService, productService service.ProductService, categoryService service.CategoryService, userService service.UserService, validator *validator.Validate) OrderHandler {
	return OrderHandler{
		OrderService:    orderService,
		ProductService:  productService,
		CategoryService: categoryService,
		UserService:     userService,
		Validator:       validator,
	}
}

// Get godoc
//
//	@Tags			order
//	@Summary		Show a order
//	@Description	get order by ID
//	@Produce		json
//	@Param			id	path		int	true	"Order ID"
//	@Success		200	{object}	util.ApiResponse{data=model.Order}
//	@Failure		400	{object}	util.ApiResponse{}
//	@Failure		500	{object}	util.ApiResponse{}
//	@Router			/order/{id} [get]
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	var response util.ApiResponse
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid order id"
		util.WriteJson(w, response)
		return
	}

	order, err := h.OrderService.Get(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status = http.StatusBadRequest
			response.Message = "Order not found"
			util.WriteJson(w, response)
			return
		}
		response.Status = http.StatusInternalServerError
		response.Message = "Error while getting order"
		util.WriteJson(w, response)
		return
	}
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = order
	fmt.Println(order.Products)
	util.WriteJson(w, response)
}

// Create godoc
//
//	@Tags			order
//	@Summary		Create a order
//	@Description	Create a order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			order	body		dto.CreateOrderDto	true	"Order"
//	@Success		200		{object}	util.ApiResponse{}
//	@Failure		400		{object}	util.ApiResponse{}
//	@Failure		500		{object}	util.ApiResponse{}
//	@Router			/order [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data dto.CreateOrderDto
	var response util.ApiResponse
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = util.JsonDecodeError.Error()
		util.WriteJson(w, response)
		return
	}
	err := h.Validator.Struct(data)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = util.JsonDecodeError.Error()
		util.WriteJson(w, response)
		return
	}

	userID, _ := r.Context().Value(middleware.AuthUserID).(string)

	var totalAmount float64
	var orderItems []model.OrderItem

	for _, orderItem := range data.Products {
		productID := strconv.Itoa(int(orderItem.ProductID))
		product, err := h.ProductService.Get(productID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Status = http.StatusBadRequest
				response.Message = fmt.Sprintf("Product with id %d not found", orderItem.ProductID)
				util.WriteJson(w, response)
				return
			}
			response.Status = http.StatusInternalServerError
			response.Message = "Error while getting product"
			util.WriteJson(w, response)
			return
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID: orderItem.ProductID,
			Product:   product,
			Quantity:  int(orderItem.Quantity),
		})
		totalAmount += product.Price * float64(orderItem.Quantity)

	}

	userIdint, _ := strconv.Atoi(userID)
	user, _ := h.UserService.Get(userID)
	order := model.Order{
		UserID:      uint(userIdint),
		User:        user,
		TotalAmount: totalAmount,
		Products:    orderItems,
	}

	err = h.OrderService.Create(order)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while creating order"
		util.WriteJson(w, response)

	}
}
