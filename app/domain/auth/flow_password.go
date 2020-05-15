package auth

import (
	"fmt"
	"github.com/pestanko/gouthy/app/domain/users"
)

type UcPasswordFlow struct {
	Users users.Repository
}

func (flow *UcPasswordFlow) Check(user *users.User, password string) error {
	// check password
	if user.CheckPassword(password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
