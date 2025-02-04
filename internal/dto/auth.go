package dto

type Login struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type Register struct {
	Name     string `validate:"required"`
	LastName string `validate:"required"`
	UserName string `validate:"required"`
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
