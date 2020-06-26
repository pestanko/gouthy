package users

import "context"

type GetUsersService interface {
	Get(ctx context.Context, params UserQuery) ([]User, error)
	GetOne(ctx context.Context, params UserQuery) (*User, error)
}

type GetUsersServiceImpl struct {
	repo Repository
}

func (s *GetUsersServiceImpl) Get(ctx context.Context, params UserQuery) (result []User, err error) {
	return s.repo.Query(ctx, params)
}

func (s *GetUsersServiceImpl) GetOne(ctx context.Context, params UserQuery) (*User, error) {
	return s.repo.QueryOne(ctx, params)
}

func NewGetUsersService(usersRepo Repository) GetUsersService {
	return &GetUsersServiceImpl{repo: usersRepo}
}
