package jwt

import (
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type jwtRequest struct {
	clientID  string `json:"client_id" bson:"client_id"`
	userEmail string `json:"email" bson:"email"`
	expire    int64  `json:"expire" bson:"expire"`
	grantType string `json:"grant_type" bson:"grant_type"`
}

func NewJwtRequest(clientID string, grantType string, expire int64) *jwtRequest {
	return &jwtRequest{
		clientID:  clientID,
		expire:    expire,
		grantType: clientID,
	}
}

func (j *jwtRequest) GenerateToken(tokenType string) (string, error) {
	var secretKey string

	if tokenType == "access_token" {
		secretKey = os.Getenv("JWT_ACCESS_SECRET_KEY")
	} else if tokenType == "login_token" {
		secretKey = os.Getenv("JWT_LOGIN_SECRET_KEY")
	} else {
		secretKey = os.Getenv("JWT_REFRESH_SECRET_KEY")
	}

	atClaims := jwt.MapClaims{}
	atClaims["client_id"] = j.clientID
	atClaims["grant_type"] = j.grantType
	atClaims["expire"] = j.expire
	atClaims["email"] = j.userEmail
	atClaims["token_uuid"] = uuid.New().String()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func ExtractToken(token string) string {
	strArr := strings.Split(token, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(tokenString string, tokenType string) (*jwt.Token, error) {
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

func IsValidToken(token *jwt.Token) error {
	var err error

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func TokenDecode(token string) (*models.UserTokenBodyJwt, error) {
	tokenJwt, err := jwt.ParseWithClaims(token, &models.UserTokenBodyJwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_LOGIN_SECRET_KEY")), nil
	})
	if err != nil {
		errString := fmt.Sprintf("Error parse token, %s", err.Error())
		err := helper.NewError(errString)
		return nil, err
	}

	claimsToken, _ := tokenJwt.Claims.(*models.UserTokenBodyJwt)

	return claimsToken, nil
}
