package apps

import "github.com/pestanko/gouthy/app/shared"

func NewApplicationError(msg string) shared.AppError {
	return shared.NewAppError(msg).WithType("application")
}