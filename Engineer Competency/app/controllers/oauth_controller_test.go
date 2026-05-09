package controllers

import (
	"bytes"
	"encoding/json"
	"engineer-comp/app/global/utils/model"
	"engineer-comp/app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOauthController_GenerateToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOauthUseCase := new(MockOauthUseCase)
	controller := InitOauthController(mockOauthUseCase)

	request := models.Request{
		ClientID:     "test_client",
		ClientSecret: "test_secret",
		Username:     "test@example.com",
		Password:     "password",
	}

	expectedResponse := &models.TokenResponse{
		AccessToken: "access_token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	mockOauthUseCase.On("GenerateToken", mock.AnythingOfType("*models.Request")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GenerateToken(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockOauthUseCase.AssertExpectations(t)
}

func TestOauthController_GenerateToken_MissingClientID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOauthUseCase := new(MockOauthUseCase)
	controller := InitOauthController(mockOauthUseCase)

	request := models.Request{
		ClientSecret: "test_secret",
		Username:     "test@example.com",
		Password:     "password",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GenerateToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Error, "Client ID is Required")
}

func TestOauthController_VerifyAndValidateLoginToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOauthUseCase := new(MockOauthUseCase)
	controller := InitOauthController(mockOauthUseCase)

	request := models.Request{
		UserToken: "valid_token",
	}

	expectedResponse := &models.ValidLoginTokenWithClient{
		IsValid: true,
		TokenBody: &models.UserTokenBodyJwt{
			ClientId: "test_client",
			Email:    "test@example.com",
		},
	}

	mockOauthUseCase.On("VerifyAndValidateLoginToken", "valid_token").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/oauth/verify-login-token", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.VerifyAndValidateLoginToken(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockOauthUseCase.AssertExpectations(t)
}

func TestOauthController_VerifyAndValidateLoginToken_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOauthUseCase := new(MockOauthUseCase)
	controller := InitOauthController(mockOauthUseCase)

	request := models.Request{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/oauth/verify-login-token", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.VerifyAndValidateLoginToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Error, "Login Token is Required")
}
