package apps

import "github.com/pestanko/gouthy/app/shared"

func NewApplicationError(msg string) shared.GouthyError {
	return shared.NewGouthyError(msg).WithType("application")
}