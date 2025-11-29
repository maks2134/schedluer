package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"schedluer/internal/models"
)

type FavoriteRepository interface {
	GetAll(ctx context.Context, userID string) ([]models.FavoriteGroup, error)
	GetByGroupNumber(ctx context.Context, userID string, groupNumber string) (*models.FavoriteGroup, error)
	Search(ctx context.Context, userID string, query string) ([]models.FavoriteGroup, error)
	Add(ctx context.Context, favorite *models.FavoriteGroup) error
	Delete(ctx context.Context, userID string, groupNumber string) error
	IsFavorite(ctx context.Context, userID string, groupNumber string) (bool, error)
}

type favoriteRepository struct {
	collection *mongo.Collection
	logger     *logrus.Logger
}

func NewFavoriteRepository(db *mongo.Database, logger *logrus.Logger) FavoriteRepository {
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
		logger:     logger,
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
		var rawDoc bson.M
		if err := cursor.Decode(&rawDoc); err != nil {
			if r.logger != nil {
				r.logger.Warnf("Failed to decode favorite document (raw): %v", err)
			}
			continue
		}

		favorite := models.FavoriteGroup{}

		if id, ok := rawDoc["_id"]; ok {
			if idArray, ok := id.([]interface{}); ok && len(idArray) > 0 {
				id = idArray[0]
				if r.logger != nil {
					r.logger.Debugf("Fixed _id from array to single value")
				}
			} else if idArray, ok := id.(bson.A); ok && len(idArray) > 0 {
				id = idArray[0]
				if r.logger != nil {
					r.logger.Debugf("Fixed _id from bson.A to single value")
				}
			}

			if objID, ok := id.(primitive.ObjectID); ok {
				favorite.ID = objID
			} else if idStr, ok := id.(string); ok {
				if objID, err := primitive.ObjectIDFromHex(idStr); err == nil {
					favorite.ID = objID
				}
			}
		}

		if groupNumber, ok := rawDoc["group_number"].(string); ok {
			favorite.GroupNumber = groupNumber
		}
		if userID, ok := rawDoc["user_id"].(string); ok {
			favorite.UserID = userID
		}
		if createdAt, ok := rawDoc["created_at"].(primitive.DateTime); ok {
			favorite.CreatedAt = createdAt.Time()
		} else if createdAt, ok := rawDoc["created_at"].(time.Time); ok {
			favorite.CreatedAt = createdAt
		}
		if updatedAt, ok := rawDoc["updated_at"].(primitive.DateTime); ok {
			favorite.UpdatedAt = updatedAt.Time()
		} else if updatedAt, ok := rawDoc["updated_at"].(time.Time); ok {
			favorite.UpdatedAt = updatedAt
		}

		favorites = append(favorites, favorite)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if favorites == nil {
		favorites = []models.FavoriteGroup{}
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

func (r *favoriteRepository) Search(ctx context.Context, userID string, query string) ([]models.FavoriteGroup, error) {
	pipeline := []bson.M{
		{
			"$search": bson.M{
				"index": "default",
				"compound": bson.M{
					"must": []bson.M{
						{
							"text": bson.M{
								"query": query,
								"path":  "group_number",
							},
						},
						{
							"equals": bson.M{
								"value": userID,
								"path":  "user_id",
							},
						},
					},
				},
			},
		},
		{
			"$limit": 100,
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		if r.logger != nil {
			r.logger.Warnf("Search index not available, falling back to regex search: %v", err)
		}
		filter := bson.M{
			"user_id":      userID,
			"group_number": bson.M{"$regex": query, "$options": "i"},
		}
		cursor, err = r.collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var favorites []models.FavoriteGroup
	for cursor.Next(ctx) {
		var rawDoc bson.M
		if err := cursor.Decode(&rawDoc); err != nil {
			if r.logger != nil {
				r.logger.Warnf("Failed to decode favorite document (raw) in search: %v", err)
			}
			continue
		}

		// Создаем структуру вручную, чтобы избежать проблем с _id
		favorite := models.FavoriteGroup{}

		// Обрабатываем _id отдельно - если это массив, берем первый элемент
		if id, ok := rawDoc["_id"]; ok {
			if idArray, ok := id.([]interface{}); ok && len(idArray) > 0 {
				id = idArray[0]
			} else if idArray, ok := id.(bson.A); ok && len(idArray) > 0 {
				id = idArray[0]
			}

			if objID, ok := id.(primitive.ObjectID); ok {
				favorite.ID = objID
			} else if idStr, ok := id.(string); ok {
				if objID, err := primitive.ObjectIDFromHex(idStr); err == nil {
					favorite.ID = objID
				}
			}
		}

		if groupNumber, ok := rawDoc["group_number"].(string); ok {
			favorite.GroupNumber = groupNumber
		}
		if userID, ok := rawDoc["user_id"].(string); ok {
			favorite.UserID = userID
		}
		if createdAt, ok := rawDoc["created_at"].(primitive.DateTime); ok {
			favorite.CreatedAt = createdAt.Time()
		} else if createdAt, ok := rawDoc["created_at"].(time.Time); ok {
			favorite.CreatedAt = createdAt
		}
		if updatedAt, ok := rawDoc["updated_at"].(primitive.DateTime); ok {
			favorite.UpdatedAt = updatedAt.Time()
		} else if updatedAt, ok := rawDoc["updated_at"].(time.Time); ok {
			favorite.UpdatedAt = updatedAt
		}

		favorites = append(favorites, favorite)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if favorites == nil {
		favorites = []models.FavoriteGroup{}
	}

	return favorites, nil
}
