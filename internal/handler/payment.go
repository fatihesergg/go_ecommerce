package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/fatihesergg/go_ecommerce/internal/model"
	"github.com/fatihesergg/go_ecommerce/internal/service"
	"github.com/fatihesergg/go_ecommerce/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v81"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	PaymentService service.PaymentService
	OrderService   service.OrderService
	Validator      *validator.Validate
}

func NewPaymentHandler(paymentService service.PaymentService, orderService service.OrderService, validator *validator.Validate) PaymentHandler {
	return PaymentHandler{PaymentService: paymentService, OrderService: orderService, Validator: validator}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	STRIPE_API := os.Getenv("STRIPE_API")

	defer r.Body.Close()
	var response util.ApiResponse

	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid order id"
		util.WriteJson(w, response)
		return
	}

	if STRIPE_API == "" {
		response.Status = http.StatusInternalServerError
		response.Message = "Something went wrong"
		util.WriteJson(w, response)
		return
	}

	stripe.Key = STRIPE_API

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

	payment := model.Payment{
		Amount:  float64(order.TotalAmount),
		OrderID: order.ID,
		Order:   order,
	}

	err = h.PaymentService.ProcessPayment(payment)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error while saving payment"
		util.WriteJson(w, response)
		return
	}

	response.Status = http.StatusOK
	response.Message = "Success"
	util.WriteJson(w, response)
}
