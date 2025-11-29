package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteGroup struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupNumber string             `bson:"group_number" json:"group_number"`
	UserID      string             `bson:"user_id" json:"user_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
