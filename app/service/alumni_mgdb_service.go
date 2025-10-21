package service

import (
	"context"
	"go_clean/app/models"
	"go_clean/app/repository"
	"time"
)

type AlumniMongoService struct {
	repo *repository.AlumniMongoRepository
}

func NewAlumniMongoService(repo *repository.AlumniMongoRepository) *AlumniMongoService {
	return &AlumniMongoService{repo: repo}
}

func (s *AlumniMongoService) Create(ctx context.Context, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return s.repo.Create(ctx, data)
}

func (s *AlumniMongoService) GetAll(ctx context.Context) ([]models.AlumniMongo, error) {
	return s.repo.FindAll(ctx)
}

func (s *AlumniMongoService) GetByID(ctx context.Context, id string) (*models.AlumniMongo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AlumniMongoService) Update(ctx context.Context, id string, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	data.UpdatedAt = time.Now()
	return s.repo.Update(ctx, id, data)
}

func (s *AlumniMongoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
