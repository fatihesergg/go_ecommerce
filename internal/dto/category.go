package dto

type CategoryCreateDto struct {
	Name string `json:"name" validate:"required"`
}

type CategoryUpdateDto struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
