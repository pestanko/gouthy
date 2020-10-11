package users

import (
	"github.com/pestanko/gouthy/app/shared"
)

func NewErrInvalidPassword() shared.AppError {
	return shared.NewErrInvalidField("password")
}

func NewErrInvalidUsername() shared.AppError {
	return shared.NewErrInvalidField("username")
}


func NewErrInvalidEmail() shared.AppError {
	return shared.NewErrInvalidField("email")
}


