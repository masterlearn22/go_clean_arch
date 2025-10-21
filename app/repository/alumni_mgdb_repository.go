package repository

import (
	"context"
	"fmt"
	"go_clean/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlumniMongoRepository struct {
	collection *mongo.Collection
}

func NewAlumniMongoRepository(db *mongo.Database) *AlumniMongoRepository {
	return &AlumniMongoRepository{
		collection: db.Collection("alumni"),
	}
}

// Create
func (r *AlumniMongoRepository) Create(ctx context.Context, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	result, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = result.InsertedID.(primitive.ObjectID)
	return data, nil
}

// Get All
func (r *AlumniMongoRepository) FindAll(ctx context.Context) ([]models.AlumniMongo, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var alumniList []models.AlumniMongo
	if err := cur.All(ctx, &alumniList); err != nil {
		return nil, err
	}
	return alumniList, nil
}

// Get By ID (bisa _id Mongo atau alumni_id custom)
func (r *AlumniMongoRepository) FindByID(ctx context.Context, id string) (*models.AlumniMongo, error) {
	var result models.AlumniMongo

	// coba dulu convert ke ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
		if err == nil {
			return &result, nil
		}
	}

	// fallback ke alumni_id numerik
	err = r.collection.FindOne(ctx, bson.M{"alumni_id": id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("data tidak ditemukan")
	}
	return &result, nil
}

// Update
func (r *AlumniMongoRepository) Update(ctx context.Context, id string, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	filter := bson.M{}

	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{"alumni_id": id}
	}

	update := bson.M{"$set": data}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Delete
func (r *AlumniMongoRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{}

	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{"alumni_id": id}
	}

	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
