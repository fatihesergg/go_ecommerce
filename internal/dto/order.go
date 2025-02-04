package dto

import (
	"database/sql"
)

type OrderProductDto struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `json:"name" `
	ImageURL   sql.NullString `json:"image_url"`
	Price      float64        `json:"price" `
	Stock      uint           `json:"stock" `
	CategoryID uint           `json:"category_id" `
}

type OrderItemDto struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity"  validate:"required"`
}
type CreateOrderDto struct {
	Products []OrderItemDto `json:"products" validate:"required" `
}

type UpdateOrderDto struct {
	ID          int            `json:"id" validate:"required"`
	Products    []OrderItemDto `json:"products" validate:"required"`
	TotalAmount float64        `json:"total_amount" validate:"required"`
}
