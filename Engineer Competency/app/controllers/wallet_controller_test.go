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

func TestWalletController_CreateWallet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	request := models.Wallet{
		UserID:       1,
		CurrentSaldo: 0,
	}

	expectedWallet := &models.Wallet{
		ID:           1,
		UserID:       1,
		CurrentSaldo: 0,
	}

	mockWalletUseCase.On("CreateWallet", mock.AnythingOfType("*models.Wallet")).Return(expectedWallet, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/wallet/create", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateWallet(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockWalletUseCase.AssertExpectations(t)
}

func TestWalletController_TopUpWallet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	request := models.Wallet{
		UserID:       1,
		CurrentSaldo: 500,
	}

	expectedWallet := &models.Wallet{
		ID:           1,
		UserID:       1,
		CurrentSaldo: 1500,
	}

	mockWalletUseCase.On("TopUpWallet", mock.AnythingOfType("*models.Wallet")).Return(expectedWallet, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/wallet/top-up", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.TopUpWallet(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockWalletUseCase.AssertExpectations(t)
}

func TestWalletController_GetWalletByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	expectedWallet := &models.Wallet{
		ID:           1,
		UserID:       1,
		CurrentSaldo: 1000,
	}

	mockWalletUseCase.On("GetWalletByID", 1).Return(expectedWallet, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("GET", "/wallet/1", nil)

	controller.GetWalletByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockWalletUseCase.AssertExpectations(t)
}

func TestWalletController_GetWalletByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	c.Request, _ = http.NewRequest("GET", "/wallet/invalid", nil)

	controller.GetWalletByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "invalid wallet id", response.Error)
}

func TestWalletController_GetWalletByUserID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	expectedWallets := []models.Wallet{
		{ID: 1, UserID: 1, CurrentSaldo: 1000},
	}

	mockWalletUseCase.On("GetWalletByUserID", 1).Return(expectedWallets, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}

	c.Request, _ = http.NewRequest("GET", "/wallet/user/1", nil)

	controller.GetWalletByUserID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockWalletUseCase.AssertExpectations(t)
}

func TestWalletController_GetWalletBalanceByUserID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockWalletUseCase := new(MockWalletUseCase)
	controller := InitWalletController(mockWalletUseCase)

	expectedBalance := 1500.0

	mockWalletUseCase.On("GetBalanceWalletByUserID", 1).Return(expectedBalance, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}

	c.Request, _ = http.NewRequest("GET", "/wallet/balance/1", nil)

	controller.GetWalletBalanceByUserID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockWalletUseCase.AssertExpectations(t)
}
