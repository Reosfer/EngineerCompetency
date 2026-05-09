package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientGrantSession struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Client    *Client            `json:"client" bson:"client"`
	Users     *LoginResponse     `json:"user_login,omitempty" bson:"users_login,omitempty"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type ClientGrantSessionChan struct {
	ClientGrantSession *ClientGrantSession
	Error              error
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type Client struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID     string             `json:"client_id" bson:"client_id"`
	ClientSecret string             `json:"client_secret" bson:"client_secret"`
	ClientName   string             `json:"client_name,omitempty" bson:"client_name,omitempty"`
}
