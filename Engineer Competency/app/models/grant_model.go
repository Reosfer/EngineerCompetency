package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Grant struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type GrantChan struct {
	Grant *Grant
	Error error
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}
