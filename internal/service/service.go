package service

import (
	"context"
	"l0/internal/domain"
	"l0/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, order domain.Order) error {
	return s.repo.Create(ctx, order)
}
