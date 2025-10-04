package service

import (
	"context"

	"telecomx-portability-service/internal/domain/model"
	"telecomx-portability-service/internal/infrastructure/adapter/repository"
)

type PortabilityService struct {
	repo *repository.MongoRepository
}

func NewPortabilityService(repo *repository.MongoRepository) *PortabilityService {
	return &PortabilityService{repo: repo}
}

func (s *PortabilityService) Create(ctx context.Context, p *model.Portability) error {
	return s.repo.Create(ctx, p)
}

func (s *PortabilityService) UpdateStatus(ctx context.Context, userID, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}

func (s *PortabilityService) Delete(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(ctx, userID)
}

func (s *PortabilityService) GetAll(ctx context.Context) ([]model.Portability, error) {
	return s.repo.GetAll(ctx)
}
