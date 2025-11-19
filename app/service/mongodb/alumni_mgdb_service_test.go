package service

import (
	"context"
	"testing"

	"go_clean/app/models/mongodb"
	"go_clean/app/repository/mongodb"
) 

func TestGetByID(t *testing.T) {
	mockRepo := repository.NewMockAlumniMongoRepository()
	svc := NewAlumniMongoService(mockRepo)
	ctx := context.Background()

	// insert data dulu
	data := &models.AlumniMongo{Nama: "Test"}
	created, _ := svc.Create(ctx, data)

	result, err := svc.GetByID(ctx, created.ID.Hex())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.Nama != "Test" {
		t.Errorf("expected Test, got %v", result.Nama)
	}
}

func TestDeleteAlumni(t *testing.T) {
	mockRepo := repository.NewMockAlumniMongoRepository()
	svc := NewAlumniMongoService(mockRepo)
	ctx := context.Background()

	data := &models.AlumniMongo{Nama: "Delete Test"}
	created, _ := svc.Create(ctx, data)

	err := svc.Delete(ctx, created.ID.Hex())
	if err != nil {
		t.Errorf("unexpected error on delete: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID.Hex())
	if err == nil {
		t.Errorf("expected error after delete, got nil")
	}
}
