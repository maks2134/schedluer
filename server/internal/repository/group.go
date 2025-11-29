package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"schedluer/internal/models"
)

type GroupRepository interface {
	GetByNumber(ctx context.Context, groupNumber string) (*models.StoredGroup, error)
	GetByID(ctx context.Context, id int) (*models.StoredGroup, error)
	GetAll(ctx context.Context) ([]models.StoredGroup, error)
	Save(ctx context.Context, group *models.StoredGroup) error
	SaveMany(ctx context.Context, groups []models.StoredGroup) error
	Update(ctx context.Context, group *models.StoredGroup) error
	Delete(ctx context.Context, id int) error
}

type groupRepository struct {
	collection *mongo.Collection
}

func NewGroupRepository(db *mongo.Database) GroupRepository {
	collection := db.Collection("groups")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "bsuir_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "group_data.name", Value: 1}},
		},
	}

	_, _ = collection.Indexes().CreateMany(context.Background(), indexes)

	return &groupRepository{
		collection: collection,
	}
}

func (r *groupRepository) GetByNumber(ctx context.Context, groupNumber string) (*models.StoredGroup, error) {
	var group models.StoredGroup
	err := r.collection.FindOne(ctx, bson.M{"group_data.name": groupNumber}).Decode(&group)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) GetByID(ctx context.Context, id int) (*models.StoredGroup, error) {
	var group models.StoredGroup
	err := r.collection.FindOne(ctx, bson.M{"bsuir_id": id}).Decode(&group)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) GetAll(ctx context.Context) ([]models.StoredGroup, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var groups []models.StoredGroup
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepository) Save(ctx context.Context, group *models.StoredGroup) error {
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *groupRepository) SaveMany(ctx context.Context, groups []models.StoredGroup) error {
	if len(groups) == 0 {
		return nil
	}

	docs := make([]interface{}, len(groups))
	for i := range groups {
		docs[i] = groups[i]
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := r.collection.InsertMany(ctx, docs, opts)
	return err
}

func (r *groupRepository) Update(ctx context.Context, group *models.StoredGroup) error {
	filter := bson.M{"bsuir_id": group.BSUIRID}

	update := bson.M{
		"$set": bson.M{
			"bsuir_id":         group.BSUIRID,
			"group_data":       group.GroupData,
			"last_update_date": group.LastUpdateDate,
			"updated_at":       group.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": group.CreatedAt,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *groupRepository) Delete(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"bsuir_id": id})
	return err
}
