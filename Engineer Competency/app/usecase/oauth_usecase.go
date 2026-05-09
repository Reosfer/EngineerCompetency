package usecases

import (
	"engineer-comp/app/global/jwt"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"engineer-comp/app/repositories"

	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	goJWT "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type OauthUseCaseInterface interface {
	GenerateToken(request *models.Request) (*models.TokenResponse, error)
	VerifyAndValidateLoginToken(tokenString string) (*models.ValidLoginTokenWithClient, error)
	ValidateTokenUser(token string, requestID string, resultChan chan *models.ValidUserTokenWithClient)
	CheckUserRole(userEmail string, roleId string) (int64, int64, error)
}

type oauthUseCase struct {
	grantRepository repositories.GrantRepositoryInterface
}

func InitOauthUseCase(grantRepository repositories.GrantRepositoryInterface) OauthUseCaseInterface {
	return &oauthUseCase{
		grantRepository: grantRepository,
	}
}

func (u *oauthUseCase) GenerateToken(request *models.Request) (*models.TokenResponse, error) {
	response := &models.TokenResponse{}
	res, err := u.LoginCredential(request)

	if err != nil {
		fmt.Println("login_internal case error, " + err.Error())
		return nil, err
	}

	response = res
	return response, nil
}

func (u *oauthUseCase) VerifyAndValidateLoginToken(tokenString string) (*models.ValidLoginTokenWithClient, error) {
	token, err := jwt.VerifyToken(tokenString, "login_token")

	if err != nil {
		fmt.Println("error verify")
		return &models.ValidLoginTokenWithClient{}, err
	}

	tokenJwt, err := goJWT.ParseWithClaims(tokenString, &models.UserTokenBodyJwt{}, func(token *goJWT.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_LOGIN_SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Println("error ParseWithClaims: ", err)
		return &models.ValidLoginTokenWithClient{}, err
	}

	claimsToken, _ := tokenJwt.Claims.(*models.UserTokenBodyJwt)

	resultGetTokenChan := make(chan *models.LoginResponseChan)
	go u.grantRepository.GetLoginSession(tokenString, resultGetTokenChan)
	resultGetToken := <-resultGetTokenChan

	if resultGetToken.Error != nil {
		fmt.Println("get token error, " + resultGetToken.Error.Error())
		return &models.ValidLoginTokenWithClient{}, resultGetToken.Error
	}

	now := time.Now().Unix()

	if resultGetToken.LoginToken.ExpireTime < now {
		fmt.Println("token error, token expired")
		err = helper.NewError("token expired")
		return &models.ValidLoginTokenWithClient{}, err
	}

	err = jwt.IsValidToken(token)

	if err != nil {
		fmt.Println("jwt error, " + err.Error())
		return &models.ValidLoginTokenWithClient{}, err
	}

	response := &models.ValidLoginTokenWithClient{
		IsValid: true,
		Client: &models.Client{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
		TokenBody: claimsToken,
	}

	return response, nil
}

func (u *oauthUseCase) Login(request *models.Login) (*models.LoginResponse, error) {
	resultPassword := &models.LoginResponseChan{}
	getLoginChan := make(chan *models.LoginResponseChan)

	//param check email or nik
	if ValidEmail(request.UserName) {
		//email
		go u.grantRepository.GetPasswordByEmail(request.UserName, getLoginChan)
		resultPassword = <-getLoginChan
	} else {
		fmt.Println("Wrong Email Format.")
		err := helper.NewError("Wrong Email Format")
		return nil, err
	}

	if resultPassword.Error != nil {
		return nil, resultPassword.Error
	}
	fmt.Println("check password")
	match := CheckPasswordHash(request.Password, resultPassword.Login.Password)

	if !match {
		fmt.Println("Password Is Wrong.")
		err := helper.NewError("Wrong password")
		return nil, err
	}

	var loginRes *models.LoginResponse
	loginRes = resultPassword.Login

	return loginRes, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *oauthUseCase) LoginCredential(request *models.Request) (*models.TokenResponse, error) {
	//checking the users login first
	fmt.Println("HIT LOGIN CREDENTIAL V1")
	var login models.Login
	login.UserName = request.Username
	login.Password = request.Password
	loginRes, err := u.Login(&login)
	if err != nil {
		return &models.TokenResponse{}, err
	}

	fmt.Println("HIT LOGIN CREDENTIAL V1 - 2")
	login.Password = "********************"
	loginRes.Password = login.Password

	//starting generate token login
	clientID := strings.ToLower(request.ClientID)
	clientSecret := strings.ToLower(request.ClientSecret)

	if clientID != strings.ToLower(os.Getenv("CLIENT_ID")) || clientSecret != strings.ToLower(os.Getenv("CLIENT_SECRET")) {
		err := helper.NewError("Invalid client credentials")
		return &models.TokenResponse{}, err
	}

	expire := time.Now().Add(time.Hour * 1).Unix()
	// newToken := jwt.NewJwtUserRequest(clientID, loginRes.UserID, loginRes.UserName, loginRes.UserNik, loginRes.Email, int64(loginRes.PositionID), expire)
	newToken := jwt.NewJwtUserRequest(clientID, clientSecret, loginRes.UserID, loginRes.UserName, loginRes.Email, "login credentials", expire)
	token, err := newToken.GenerateToken("login_token")

	if err != nil {
		fmt.Println("new token generate error, " + err.Error())
		return &models.TokenResponse{}, err
	}

	now := time.Now()
	loginTokenRequest := &models.LoginToken{
		Token:      token,
		ExpireTime: expire,
		ExpireIn:   9999,
		CreatedAt:  &now,
	}

	resultInsertSessionLoginChan := make(chan *models.LoginResponseChan)
	go u.grantRepository.InsertLoginSession(loginTokenRequest, resultInsertSessionLoginChan)
	resultInsertSessionLogin := <-resultInsertSessionLoginChan

	if resultInsertSessionLogin.Error != nil {
		fmt.Println("login token error, " + resultInsertSessionLogin.Error.Error())
		return nil, resultInsertSessionLogin.Error
	}

	tokenResponse := &models.TokenResponse{
		UsersToken: token,
		ExpiresIn:  9999,
		TokenType:  "Bearer",
	}

	return tokenResponse, nil
}

func (u *oauthUseCase) ValidateTokenUser(token string, requestID string, resultChan chan *models.ValidUserTokenWithClient) {
	response := &models.ValidUserTokenWithClient{}
	_, err := u.VerifyAndValidateLoginToken(token)

	if err != nil {
		if strings.Contains(err.Error(), "token expired") {
			response.StatusCode = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			response.StatusCode = http.StatusNotFound
		} else {
			response.StatusCode = http.StatusInternalServerError
		}
		response.Error = err
		resultChan <- response
		return
	}

	response.StatusCode = http.StatusOK
	response.Error = nil
	resultChan <- response
	fmt.Println("success validate users token")
	return

}

func (u *oauthUseCase) CheckUserRole(userEmail string, roleId string) (int64, int64, error) {

	roleIdInt, userId, err := u.grantRepository.GetRoleById(userEmail)

	if err != nil {
		fmt.Println("get role by id error, " + err.Error())
		return 0, 0, err
	}

	for _, v := range roleIdInt {
		if fmt.Sprint(v) == roleId {
			return v, userId, nil
		}
	}

	fmt.Println("user not have access to this resource")
	err = helper.NewError("user not have access to this resource")
	return 0, 0, err
}
