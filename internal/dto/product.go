package dto

type ProductCreateDto struct {
	Name       string  `json:"name" validate:"required"`
	ImageURL   *string `json:"image_url" `
	Price      float64 `json:"price" validate:"required"`
	Stock      uint    `json:"stock" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

type ProductUpdateDto struct {
	ID         int     `json:"id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	ImageURL   *string `json:"image_url"`
	Price      float64 `json:"price" validate:"required"`
	Stock      uint    `json:"stock" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}
