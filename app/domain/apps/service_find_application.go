package apps

import (
	"context"
)

type FindService interface {
	Find(ctx context.Context, params FindQuery) ([]Application, error)
	FindOne(ctx context.Context, params FindQuery) (*Application, error)
}

func NewFindAppsService(repo Repository) FindService {
	return &findApplicationsServiceImpl{repo: repo}
}

type findApplicationsServiceImpl struct {
	repo Repository
}

func (s *findApplicationsServiceImpl) Find(ctx context.Context, params FindQuery) ([]Application, error) {
	return s.repo.Query(ctx, params)
}

func (s *findApplicationsServiceImpl) FindOne(ctx context.Context, params FindQuery) (*Application, error) {
	return s.repo.QueryOne(ctx, params)
}
