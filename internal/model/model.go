package model

import (
	"time"
)

type Status int

const (
	_              = iota
	PENDING Status = iota
	SUCCESS
	FAILED
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" `
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" `
	LastName  string    `json:"last_name"  `
	UserName  string    `json:"user_name" `
	Email     string    `json:"email" `
	Role      string    `json:"role"`
	Password  string    `json:"-" `
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Comment   string    `json:"comment" `
	ProductID uint      `json:"product_id" `
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	UserID    uint      `json:"user_id" `
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Product struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name" `
	ImageURL   *string   `json:"image_url"`
	Price      float64   `json:"price" `
	Stock      uint      `json:"stock" `
	CategoryID uint      `json:"category_id" `
	Category   Category  `gorm:"foreignKey:CategoryID" json:"-"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ProductID uint    `json:"product_id" `
	Product   Product `gorm:"foreignKey:ProductID"  json:"-"`
	Quantity  int     `json:"quantity"`
	OrderID   int
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Payment struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionId string  `json:"transaction_id" `
	Amount        float64 `json:"amount" `
	Status        Status  `json:"status" `
	OrderID       uint
	Order         Order     `gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `json:"user_id"  `
	User        User        `gorm:"foreignKey:UserID" json:"-"`
	Products    []OrderItem `json:"products" gorm:"foreignKey:OrderID" `
	TotalAmount float64     `json:"total_amount" `
}
