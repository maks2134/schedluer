package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"schedluer/internal/models"
)

type FavoriteRepository interface {
	GetAll(ctx context.Context, userID string) ([]models.FavoriteGroup, error)
	GetByGroupNumber(ctx context.Context, userID string, groupNumber string) (*models.FavoriteGroup, error)
	Add(ctx context.Context, favorite *models.FavoriteGroup) error
	Delete(ctx context.Context, userID string, groupNumber string) error
	IsFavorite(ctx context.Context, userID string, groupNumber string) (bool, error)
}

type favoriteRepository struct {
	collection *mongo.Collection
}

func NewFavoriteRepository(db *mongo.Database) FavoriteRepository {
	collection := db.Collection("favorite_groups")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "group_number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	}

	_, _ = collection.Indexes().CreateMany(context.Background(), indexes)

	return &favoriteRepository{
		collection: collection,
	}
}

func (r *favoriteRepository) GetAll(ctx context.Context, userID string) ([]models.FavoriteGroup, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var favorites []models.FavoriteGroup
	for cursor.Next(ctx) {
		var favorite models.FavoriteGroup
		if err := cursor.Decode(&favorite); err != nil {
			continue
		}
		favorites = append(favorites, favorite)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
}

func (r *favoriteRepository) GetByGroupNumber(ctx context.Context, userID string, groupNumber string) (*models.FavoriteGroup, error) {
	var favorite models.FavoriteGroup
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID, "group_number": groupNumber}).Decode(&favorite)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) Add(ctx context.Context, favorite *models.FavoriteGroup) error {
	opts := options.UpdateOne().SetUpsert(true)
	filter := bson.M{"user_id": favorite.UserID, "group_number": favorite.GroupNumber}
	update := bson.M{
		"$set": bson.M{
			"group_number": favorite.GroupNumber,
			"user_id":      favorite.UserID,
			"updated_at":   favorite.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": favorite.CreatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *favoriteRepository) Delete(ctx context.Context, userID string, groupNumber string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID, "group_number": groupNumber})
	return err
}

func (r *favoriteRepository) IsFavorite(ctx context.Context, userID string, groupNumber string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID, "group_number": groupNumber})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
