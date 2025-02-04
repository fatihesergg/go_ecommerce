package dto

type ReviewCreateDto struct {
	Comment   string `json:"comment" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
}

type ReviewUpdateDto struct {
	ID        int    `json:"id" validate:"required"`
	Comment   string `json:"comment" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
}
