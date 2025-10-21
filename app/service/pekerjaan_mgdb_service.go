package service

import (
	"context"
	"go_clean/app/models"
	"go_clean/app/repository"
	"time"
)

type PekerjaanMongoService struct {
	Repo *repository.PekerjaanMongoRepository
}

func NewPekerjaanMongoService(repo *repository.PekerjaanMongoRepository) *PekerjaanMongoService {
	return &PekerjaanMongoService{Repo: repo}
}

func (s *PekerjaanMongoService) Create(ctx context.Context, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return s.Repo.Create(ctx, p)
}

func (s *PekerjaanMongoService) GetAll(ctx context.Context) ([]models.PekerjaanMongo, error) {
	return s.Repo.FindAll(ctx)
}

func (s *PekerjaanMongoService) GetByID(ctx context.Context, id string) (*models.PekerjaanMongo, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *PekerjaanMongoService) GetByAlumniID(ctx context.Context, alumniID int) ([]models.PekerjaanMongo, error) {
	return s.Repo.FindByAlumniID(ctx, alumniID)
}

func (s *PekerjaanMongoService) Update(ctx context.Context, id string, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	p.UpdatedAt = time.Now()
	return s.Repo.Update(ctx, id, p)
}

func (s *PekerjaanMongoService) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
