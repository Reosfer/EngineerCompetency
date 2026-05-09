package jwt

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type jwtUserRequest struct {
	userID       int    `json:"user_id" bson:"user_id"`
	userName     string `json:"user_name" bson:"user_name"`
	clientID     string `json:"client_id" bson:"client_id"`
	clientSecret string `json:"client_secret" bson:"client_secret"`
	userEmail    string `json:"email" bson:"user_nik"`
	expire       int64  `json:"expire" bson:"expire"`
	grantType    string `json:"grant_type" bson:"grant_type"`
}

func NewJwtUserRequest(clientID string, clientSecret string, userID int, userName string, email string, grantType string, expire int64) *jwtUserRequest {
	return &jwtUserRequest{
		userID:       userID,
		userName:     userName,
		clientID:     clientID,
		clientSecret: clientSecret,
		userEmail:    email,
		expire:       expire,
		grantType:    clientID,
	}
}

func (j *jwtUserRequest) GenerateToken(tokenType string) (string, error) {
	var secretKey string

	secretKey = os.Getenv("JWT_LOGIN_SECRET_KEY")

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = j.userID
	atClaims["user_name"] = j.userName
	atClaims["client_id"] = j.clientID
	atClaims["client_secret"] = j.clientSecret
	atClaims["email"] = j.userEmail
	atClaims["expire"] = j.expire
	atClaims["grant_type"] = j.grantType
	atClaims["token_uuid"] = uuid.New().String()

	//atClaims["user_id"] = j.userID

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func ExtractUserToken(token string) string {
	strArr := strings.Split(token, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyUserToken(tokenString string, tokenType string) (*jwt.Token, error) {
	var secretKey string

	if tokenType == "access_token" {
		secretKey = os.Getenv("JWT_ACCESS_SECRET_KEY")
	} else if tokenType == "login_token" {
		secretKey = os.Getenv("JWT_LOGIN_SECRET_KEY")
	} else {
		secretKey = os.Getenv("JWT_REFRESH_SECRET_KEY")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func IsValidUserToken(token *jwt.Token) error {
	var err error

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}
