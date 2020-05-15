package machines

import (
	uuid "github.com/satori/go.uuid"
)

type MachineDTO struct {
	ID uuid.UUID
}

type Facade interface {
	GetByCodename(codename string) (*MachineDTO, error)
}

type FacadeImpl struct {
}

func (f FacadeImpl) GetByCodename(codename string) (*MachineDTO, error) {
	panic("implement me")
}

func NewMachinesFacade(machines Repository) Facade {
	return &FacadeImpl{}
}
