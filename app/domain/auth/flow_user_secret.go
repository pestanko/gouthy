package auth

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared/utils"
	uuid "github.com/satori/go.uuid"
)

type UcSecretFlow struct {
	users       users.Repository
	userSecrets users.SecretsRepository
	ctx         context.Context
}

func (flow *UcSecretFlow) Check(codename, secret string) error {
	// Get User
	id, err := flow.getUserIdByCodename(codename)
	if err != nil {
		return err
	}

	// check secret
	secrets, err := flow.userSecrets.List(flow.ctx, id)
	if err != nil {
		return err
	}

	if !flow.checkSecret(secrets, secret) {
		return fmt.Errorf("no valid secret found")
	}

	return nil
}

func (flow *UcSecretFlow) checkSecret(secrets []users.Secret, secret string) bool {
	for _, s := range secrets {
		if utils.CompareHashAndOriginal(s.Value, secret) && !s.IsExpired() {
			return true
		}
	}
	return false
}

func (flow *UcSecretFlow) getUserIdByCodename(username string) (uuid.UUID, error) {
	user, err := flow.users.FindByUsername(flow.ctx, username)
	if err != nil {
		return uuid.UUID{}, err
	}
	return user.ID, nil
}
