package service

import (
	"context"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}
func (s *Service) GetOrder(ctx context.Context, uid string) (*domain.Order, error) {
	return s.repo.Get(ctx, uid)
}
func (s *Service) CreateOrder(ctx context.Context, order domain.Order) error {
	return s.repo.Create(ctx, order)
}
