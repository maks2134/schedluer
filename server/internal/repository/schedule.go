package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"schedluer/internal/models"
)

type ScheduleRepository interface {
	GetByGroupNumber(ctx context.Context, groupNumber string) (*models.StoredSchedule, error)
	GetByEmployeeURLID(ctx context.Context, urlID string) (*models.StoredSchedule, error)
	Save(ctx context.Context, schedule *models.StoredSchedule) error
	Update(ctx context.Context, schedule *models.StoredSchedule) error
	Delete(ctx context.Context, groupNumber string) error
	DeleteByEmployeeURLID(ctx context.Context, urlID string) error
}

type scheduleRepository struct {
	collection *mongo.Collection
}

func NewScheduleRepository(db *mongo.Database) ScheduleRepository {
	collection := db.Collection("schedules")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "group_number", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "employee_url_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}

	_, _ = collection.Indexes().CreateMany(context.Background(), indexes)

	return &scheduleRepository{
		collection: collection,
	}
}

func (r *scheduleRepository) GetByGroupNumber(ctx context.Context, groupNumber string) (*models.StoredSchedule, error) {
	var schedule models.StoredSchedule
	err := r.collection.FindOne(ctx, bson.M{"group_number": groupNumber}).Decode(&schedule)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) GetByEmployeeURLID(ctx context.Context, urlID string) (*models.StoredSchedule, error) {
	var schedule models.StoredSchedule
	err := r.collection.FindOne(ctx, bson.M{"employee_url_id": urlID}).Decode(&schedule)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Save(ctx context.Context, schedule *models.StoredSchedule) error {
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, schedule)
	return err
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *models.StoredSchedule) error {
	schedule.UpdatedAt = time.Now()

	filter := bson.M{
		"$or": []bson.M{
			{"group_number": schedule.GroupNumber},
			{"employee_url_id": schedule.EmployeeURLID},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"schedule_data":    schedule.ScheduleData,
			"last_update_date": schedule.LastUpdateDate,
			"updated_at":       schedule.UpdatedAt,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *scheduleRepository) Delete(ctx context.Context, groupNumber string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"group_number": groupNumber})
	return err
}

func (r *scheduleRepository) DeleteByEmployeeURLID(ctx context.Context, urlID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"employee_url_id": urlID})
	return err
}
