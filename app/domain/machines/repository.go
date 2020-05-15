package machines

import "github.com/jinzhu/gorm"

type Repository interface {
	GetByCodename(codename string) (Machine, error)
}

type RepositoryDB struct {

}

func (r RepositoryDB) GetByCodename(codename string) (Machine, error) {
	panic("implement me")
}

func NewMachinesRepositoryDB(*gorm.DB) Repository {
	return &RepositoryDB {}
}