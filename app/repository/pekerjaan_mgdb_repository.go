package repository

import (
	"context"
	"go_clean/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PekerjaanMongoRepository struct {
	collection *mongo.Collection
}

func NewPekerjaanMongoRepository(db *mongo.Database) *PekerjaanMongoRepository {
	return &PekerjaanMongoRepository{
		collection: db.Collection("pekerjaan"),
	}
}

func (r *PekerjaanMongoRepository) Create(ctx context.Context, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	result, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = result.InsertedID.(primitive.ObjectID)
	return p, nil
}

func (r *PekerjaanMongoRepository) FindAll(ctx context.Context) ([]models.PekerjaanMongo, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []models.PekerjaanMongo
	if err = cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanMongoRepository) FindByID(ctx context.Context, id string) (*models.PekerjaanMongo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var result models.PekerjaanMongo
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PekerjaanMongoRepository) FindByAlumniID(ctx context.Context, alumniID int) ([]models.PekerjaanMongo, error) {
	cur, err := r.collection.Find(ctx, bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []models.PekerjaanMongo
	if err = cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanMongoRepository) Update(ctx context.Context, id string, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	_, err = r.collection.UpdateByID(ctx, objID, bson.M{"$set": p})
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *PekerjaanMongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
