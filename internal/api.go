package internal

import (
	"github.com/fatihesergg/go_ecommerce/internal/service"
	"github.com/go-playground/validator/v10"
)

type EcommerceApi *Api

type Api struct {
	CategoryService service.CategoryService
	ProductService  service.ProductService
	UserService     service.UserService
	Validator       *validator.Validate
}
