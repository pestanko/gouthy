package users

import "context"

type FindService interface {
	Find(ctx context.Context, params FindQuery) ([]User, error)
	FindOne(ctx context.Context, params FindQuery) (*User, error)
}

type findServiceImpl struct {
	repo Repository
}

func (s *findServiceImpl) Find(ctx context.Context, params FindQuery) (result []User, err error) {
	return s.repo.Query(ctx, params)
}

func (s *findServiceImpl) FindOne(ctx context.Context, params FindQuery) (*User, error) {
	return s.repo.QueryOne(ctx, params)
}

func NewFindUsersService(usersRepo Repository) FindService {
	return &findServiceImpl{repo: usersRepo}
}
