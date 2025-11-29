package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"schedluer/internal/models"
)

type EmployeeRepository interface {
	GetByURLID(ctx context.Context, urlID string) (*models.StoredEmployee, error)
	GetByID(ctx context.Context, id int) (*models.StoredEmployee, error)
	GetAll(ctx context.Context) ([]models.StoredEmployee, error)
	Save(ctx context.Context, employee *models.StoredEmployee) error
	SaveMany(ctx context.Context, employees []models.StoredEmployee) error
	Update(ctx context.Context, employee *models.StoredEmployee) error
	Delete(ctx context.Context, id int) error
}

type employeeRepository struct {
	collection *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Database) EmployeeRepository {
	collection := db.Collection("employees")

	// Создаем индексы
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "bsuir_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "url_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, _ = collection.Indexes().CreateMany(context.Background(), indexes)

	return &employeeRepository{
		collection: collection,
	}
}

func (r *employeeRepository) GetByURLID(ctx context.Context, urlID string) (*models.StoredEmployee, error) {
	var employee models.StoredEmployee
	err := r.collection.FindOne(ctx, bson.M{"url_id": urlID}).Decode(&employee)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id int) (*models.StoredEmployee, error) {
	var employee models.StoredEmployee
	err := r.collection.FindOne(ctx, bson.M{"bsuir_id": id}).Decode(&employee)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetAll(ctx context.Context) ([]models.StoredEmployee, error) {
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

	var employees []models.StoredEmployee
	for cursor.Next(ctx) {
		var employee models.StoredEmployee
		if err := cursor.Decode(&employee); err != nil {
			// Пропускаем документы с ошибками декодирования
			continue
		}
		employees = append(employees, employee)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) Save(ctx context.Context, employee *models.StoredEmployee) error {
	_, err := r.collection.InsertOne(ctx, employee)
	return err
}

func (r *employeeRepository) SaveMany(ctx context.Context, employees []models.StoredEmployee) error {
	if len(employees) == 0 {
		return nil
	}

	docs := make([]interface{}, len(employees))
	for i := range employees {
		docs[i] = employees[i]
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := r.collection.InsertMany(ctx, docs, opts)
	return err
}

func (r *employeeRepository) Update(ctx context.Context, employee *models.StoredEmployee) error {
	filter := bson.M{"bsuir_id": employee.BSUIRID}

	update := bson.M{
		"$set": bson.M{
			"bsuir_id":         employee.BSUIRID,
			"url_id":           employee.URLID,
			"employee_data":    employee.EmployeeData,
			"last_update_date": employee.LastUpdateDate,
			"updated_at":       employee.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": employee.CreatedAt,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *employeeRepository) Delete(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"bsuir_id": id})
	return err
}
