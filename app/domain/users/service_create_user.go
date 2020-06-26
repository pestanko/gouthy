package users

import "context"

type CreateUserService interface {
	CreateUser(ctx context.Context, ctrDTO *User) error
}

type createUserServiceImpl struct {
	repo Repository
}

func (c *createUserServiceImpl) CreateUser(ctx context.Context, user *User) error {
	err := c.repo.Create(ctx, user)
	return err
}

func NewCreateUserService(userRepo Repository) CreateUserService {
	return &createUserServiceImpl{
		repo: userRepo,
	}
}
