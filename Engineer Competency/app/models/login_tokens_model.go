package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginToken struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Token      string             `json:"token" bson:"token"`
	ExpireIn   int64              `json:"expire_in" bson:"expire_in"`
	ExpireTime int64              `json:"expire_time" bson:"expire_time"`
	CreatedAt  *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type LoginTokenChan struct {
	LoginToken *LoginToken
	Error      error
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type ValidLoginTokenWithClient struct {
	IsValid   bool              `json:"is_valid,omitempty" bson:"is_valid,omitempty"`
	TokenBody *UserTokenBodyJwt `json:"token_body,omitempty" bson:"token_body,omitempty"`
	Client    *Client           `json:"client,omitempty" bson:"client,omitempty"`
}

type ValidUserTokenWithClient struct {
	IsValid    bool  `json:"is_valid,omitempty" bson:"is_valid,omitempty"`
	StatusCode int   `json:"status_code,omitempty" bson:"status_code,omitempty"`
	Error      error `json:"errors,omitempty" bson:"error,omitempty"`
}
