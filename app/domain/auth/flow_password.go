package auth

import (
	"fmt"
	"github.com/pestanko/gouthy/app/domain/users"
)

type UcPasswordFlow struct {
	Users users.Repository
	user *users.User
	password string
}

func (flow *UcPasswordFlow) Check() error {
	// check password
	if ! flow.user.CheckPassword(flow.password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
