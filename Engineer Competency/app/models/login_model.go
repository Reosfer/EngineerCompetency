package models

import "github.com/golang-jwt/jwt"

type Login struct {
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
}

type LoginResponse struct {
	UserID     int    `json:"user_id" bson:"user_id"`
	UserName   string `json:"user_name" bson:"user_name"`
	UserNik    string `json:"user_nik" bson:"user_nik"`
	Email      string `json:"email" bson:"email"`
	UserStatus int    `json:"user_status" bson:"user_status"`
	Password   string `json:"password" bson:"password"`
	PositionID int    `json:"position_id" bson:"position_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UsersToken   string `json:"users_token,omitempty"`
}

type LoginResponseChan struct {
	Login       *LoginResponse
	LoginToken  *LoginToken
	Error       error
	StatusCode  int
	LoginStatus int
}

type UserTokenBody struct {
	ClientId   string `json:"client_id" bson:"client_id"`
	UserNik    string `json:"user_nik" bson:"user_nik"`
	Email      string `json:"email" bson:"email"`
	PositionID int    `json:"position_id" bson:"position_id"`
	Expire     int64  `json:"expire" bson:"expire"`
	GrantType  string `json:"grant_type" bson:"grant_type"`
	TokenUuid  string `json:"token_uuid" bson:"token_uuid"`
}
type UserTokenBodyJwt struct {
	ClientId  string `json:"client_id" bson:"client_id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email" bson:"email"`
	Expire    int64  `json:"expire" bson:"expire"`
	GrantType string `json:"grant_type" bson:"grant_type"`
	TokenUuid string `json:"token_uuid" bson:"token_uuid"`
	jwt.StandardClaims
}

type UserTokenBodyJwtV2 struct {
	ClientId  string `json:"client_id" bson:"client_id"`
	Expire    int64  `json:"expire" bson:"expire"`
	GrantType string `json:"grant_type" bson:"grant_type"`
	TokenUuid string `json:"token_uuid" bson:"token_uuid"`
	jwt.StandardClaims
}

// type jwtUserRequest struct {
// 	clientID       string `json:"client_id" bson:"client_id"`
// 	userNik        string `json:"user_nik" bson:"user_nik"`
// 	userEmail      string `json:"email" bson:"user_nik"`
// 	userPositionID int64  `json:"position_id" bson:"position_id"`
// 	expire         int64  `json:"expire" bson:"expire"`
// 	grantType      string `json:"grant_type" bson:"grant_type"`
// }
