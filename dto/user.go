package dto

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=12,max=30"`
}

type UserCreateDTO struct {
	Name     string `json:"name" validate:"required,min=4,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Country  string `json:"country"`
}