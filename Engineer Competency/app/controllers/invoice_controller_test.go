package controllers

import (
	"bytes"
	"encoding/json"
	"engineer-comp/app/global/utils/model"
	"engineer-comp/app/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInvoiceController_CreateInvoice_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockInvoiceUseCase := new(MockInvoiceUseCase)
	mockOauthUseCase := new(MockOauthUseCase)

	controller := InitInvoiceController(mockInvoiceUseCase, mockOauthUseCase)

	request := models.InvoiceRequest{
		InvoiceName: "Test Invoice",
		Amount:      1000,
		RoleID:      "2",
	}

	expectedInvoice := &models.Invoice{
		ID:          1,
		InvoiceName: "Test Invoice",
		Amount:      1000,
		Status:      "pending",
	}

	mockInvoiceUseCase.On("CreateInvoice", mock.AnythingOfType("*models.InvoiceRequest")).Return(expectedInvoice, nil)
	mockOauthUseCase.On("CheckUserRole", mock.Anything, "2").Return(int64(0), int64(1), nil)

	os.Setenv("JWT_LOGIN_SECRET_KEY", "test_secret")
	jwtClaims := jwt.MapClaims{
		"email":      "test@example.com",
		"client_id":  "test_client",
		"grant_type": "password",
		"expire":     9999999999,
		"token_uuid": "test-uuid",
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, _ := tokenObj.SignedString([]byte("test_secret"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/invoice/create", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	controller.CreateInvoice(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockInvoiceUseCase.AssertExpectations(t)
	mockOauthUseCase.AssertExpectations(t)
}

func TestInvoiceController_CreateInvoice_InvalidRoleID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockInvoiceUseCase := new(MockInvoiceUseCase)
	mockOauthUseCase := new(MockOauthUseCase)

	controller := InitInvoiceController(mockInvoiceUseCase, mockOauthUseCase)

	request := models.InvoiceRequest{
		InvoiceName: "Test Invoice",
		Amount:      1000,
		RoleID:      "",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/invoice/create", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateInvoice(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "role id is required", response.Error)
}

func TestInvoiceController_GetInvoiceByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockInvoiceUseCase := new(MockInvoiceUseCase)
	mockOauthUseCase := new(MockOauthUseCase)

	controller := InitInvoiceController(mockInvoiceUseCase, mockOauthUseCase)

	expectedInvoice := &models.Invoice{
		ID:          1,
		InvoiceName: "Test Invoice",
		Amount:      1000,
		Status:      "pending",
	}

	mockInvoiceUseCase.On("GetInvoiceByID", 1).Return(expectedInvoice, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("GET", "/invoice/1", nil)

	controller.GetInvoiceByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockInvoiceUseCase.AssertExpectations(t)
}

func TestInvoiceController_GetInvoiceByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockInvoiceUseCase := new(MockInvoiceUseCase)
	mockOauthUseCase := new(MockOauthUseCase)

	controller := InitInvoiceController(mockInvoiceUseCase, mockOauthUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	c.Request, _ = http.NewRequest("GET", "/invoice/invalid", nil)

	controller.GetInvoiceByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "invalid invoice id", response.Error)
}

func TestInvoiceController_GetAllInvoice_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockInvoiceUseCase := new(MockInvoiceUseCase)
	mockOauthUseCase := new(MockOauthUseCase)

	controller := InitInvoiceController(mockInvoiceUseCase, mockOauthUseCase)

	expectedInvoices := []models.Invoice{
		{ID: 1, InvoiceName: "Invoice 1", Amount: 1000, Status: "pending"},
		{ID: 2, InvoiceName: "Invoice 2", Amount: 2000, Status: "approved"},
	}

	mockInvoiceUseCase.On("GetAllInvoice").Return(expectedInvoices, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/invoice/all", nil)

	controller.GetAllInvoice(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, response.Data)

	mockInvoiceUseCase.AssertExpectations(t)
}
