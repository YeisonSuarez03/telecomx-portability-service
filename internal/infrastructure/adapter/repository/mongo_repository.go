package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"telecomx-portability-service/internal/domain/model"
)

type MongoRepository struct {
	Collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{Collection: db.Collection("portability")}
}

func (r *MongoRepository) Create(ctx context.Context, p *model.Portability) error {
	p.RequestedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, p)
	return err
}

func (r *MongoRepository) UpdateStatus(ctx context.Context, userID, status string) error {
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": bson.M{"current_status": status}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"userId": userID})
	return err
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]model.Portability, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Portability
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
