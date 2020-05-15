package auth

import (
	"fmt"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/machines"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared/utils"
	uuid "github.com/satori/go.uuid"
)

type UcSecretFlow struct {
	Users    users.Repository
	Machines machines.Repository
	Entities entities.Repository
	Secrets  entities.SecretsRepository
}

func (flow *UcSecretFlow) Check(codename, secret, entityType string) error {
	// Get User
	id, err := flow.getEntityIdByCodename(codename, entityType)
	if err != nil {
		return err
	}

	// check secret
	secrets, err := flow.Secrets.FindSecretsForEntity(id)
	if err != nil {
		return err
	}

	if !flow.checkSecret(secrets, secret) {
		return fmt.Errorf("no valid secret found")
	}

	return nil
}

func (flow *UcSecretFlow) checkSecret(secrets []entities.Secret, secret string) bool {
	for _, s := range secrets {
		if utils.CompareHashAndOriginal(s.Value, secret) && !s.IsExpired() {
			return true
		}
	}
	return false
}

func (flow *UcSecretFlow) getEntityIdByCodename(codename string, entityType string) (uuid.UUID, error) {
	switch entityType {
	case "machine":
		return flow.getMachineIdByCodename(codename)
	case "user":
		return flow.getUserIdByCodename(codename)
	default:
		return uuid.UUID{}, fmt.Errorf("unrecognized entity type: %s", entityType)
	}
}

func (flow *UcSecretFlow) getUserIdByCodename(codename string) (uuid.UUID, error) {
	user, err := flow.Users.FindByUsername(codename)
	if err != nil {
		return uuid.UUID{}, err
	}
	return user.ID, nil
}

func (flow *UcSecretFlow) getMachineIdByCodename(codename string) (uuid.UUID, error) {
	machine, err := flow.Machines.GetByCodename(codename)
	if err != nil {
		return uuid.UUID{}, err
	}
	return machine.ID, nil
}
