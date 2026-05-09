package controllers

import (
	"engineer-comp/app/models"

	"github.com/stretchr/testify/mock"
)

// Mock untuk InvoiceUseCaseInterface
type MockInvoiceUseCase struct {
	mock.Mock
}

func (m *MockInvoiceUseCase) CreateInvoice(request *models.InvoiceRequest) (*models.Invoice, error) {
	args := m.Called(request)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (m *MockInvoiceUseCase) GetInvoiceByID(id int) (*models.Invoice, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (m *MockInvoiceUseCase) GetAllInvoice() ([]models.Invoice, error) {
	args := m.Called()
	return args.Get(0).([]models.Invoice), args.Error(1)
}

func (m *MockInvoiceUseCase) UpdateInvoice(request *models.InvoiceRequest) (*models.Invoice, error) {
	args := m.Called(request)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (m *MockInvoiceUseCase) DeleteInvoice(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock untuk WalletUseCaseInterface
type MockWalletUseCase struct {
	mock.Mock
}

func (m *MockWalletUseCase) CreateWallet(request *models.Wallet) (*models.Wallet, error) {
	args := m.Called(request)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) TopUpWallet(request *models.Wallet) (*models.Wallet, error) {
	args := m.Called(request)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) GetWalletByID(id int) (*models.Wallet, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) GetWalletByUserID(userID int) ([]models.Wallet, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Wallet), args.Error(1)
}

func (m *MockWalletUseCase) GetBalanceWalletByUserID(userID int) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

// Mock untuk OauthUseCaseInterface
type MockOauthUseCase struct {
	mock.Mock
}

func (m *MockOauthUseCase) GenerateToken(request *models.Request) (*models.TokenResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*models.TokenResponse), args.Error(1)
}

func (m *MockOauthUseCase) VerifyAndValidateLoginToken(tokenString string) (*models.ValidLoginTokenWithClient, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*models.ValidLoginTokenWithClient), args.Error(1)
}

func (m *MockOauthUseCase) ValidateTokenUser(token string, requestID string, resultChan chan *models.ValidUserTokenWithClient) {
	// Not implemented for this test
}

func (m *MockOauthUseCase) CheckUserRole(userEmail string, roleId string) (int64, int64, error) {
	args := m.Called(userEmail, roleId)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}
