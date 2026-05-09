package usecases

import (
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"engineer-comp/app/repositories"
	"strings"

	"fmt"
)

type WalletUseCaseInterface interface {
	CreateWallet(request *models.Wallet) (*models.Wallet, error)
	TopUpWallet(request *models.Wallet) (*models.Wallet, error)
	GetWalletByID(id int) (*models.Wallet, error)
	GetWalletByUserID(userID int) ([]models.Wallet, error)
	GetBalanceWalletByUserID(userID int) (float64, error)
}

type walletUseCase struct {
	walletRepository repositories.WalletRepositoryInterface
}

func InitWalletUseCase(
	walletRepository repositories.WalletRepositoryInterface,
) WalletUseCaseInterface {

	return &walletUseCase{
		walletRepository: walletRepository,
	}
}

func (u *walletUseCase) CreateWallet(request *models.Wallet) (*models.Wallet, error) {

	if request.UserID <= 0 {
		return nil, helper.NewError("invalid user id")
	}

	if request.CurrentSaldo < 0 {
		return nil, helper.NewError("current saldo cannot be negative")
	}

	if request.RecordSaldo < 0 {
		return nil, helper.NewError("record saldo cannot be negative")
	}

	if request.Action == "" {
		return nil, helper.NewError("action is required")
	}

	createWalletChan := make(chan *models.WalletResponseChan)

	go u.walletRepository.CreateWallet(request, createWalletChan)

	result := <-createWalletChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Wallet, nil
}

func (u *walletUseCase) TopUpWallet(request *models.Wallet) (*models.Wallet, error) {

	if request.UserID <= 0 {
		return nil, helper.NewError("invalid user id")
	}

	if request.RecordSaldo < 0 {
		return nil, helper.NewError("record saldo cannot be negative")
	}

	validStatus := map[string]bool{
		"top-up": true,
	}

	if !validStatus[strings.ToLower(request.Action)] {
		return nil, helper.NewError("invalid action status")
	}

	createWalletChan := make(chan *models.WalletResponseChan)

	go u.walletRepository.TopUpWallet(request, createWalletChan)

	result := <-createWalletChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Wallet, nil
}

func (u *walletUseCase) GetWalletByID(id int) (*models.Wallet, error) {

	if id <= 0 {
		return nil, helper.NewError("invalid wallet id")
	}

	getWalletChan := make(chan *models.WalletResponseChan)

	go u.walletRepository.GetWalletByID(id, getWalletChan)

	result := <-getWalletChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Wallet, nil
}

func (u *walletUseCase) GetWalletByUserID(userID int) ([]models.Wallet, error) {

	if userID <= 0 {
		return nil, helper.NewError("invalid user id")
	}

	getWalletListChan := make(chan *models.WalletListResponseChan)

	go u.walletRepository.GetWalletByUserID(userID, getWalletListChan)

	result := <-getWalletListChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Wallet, nil
}

func (u *walletUseCase) GetBalanceWalletByUserID(userID int) (float64, error) {

	if userID <= 0 {
		return 0, helper.NewError("invalid user id")
	}

	getWalletListChan := make(chan *models.WalletResponseChan)

	go u.walletRepository.GetBalanceWalletByUserID(userID, getWalletListChan)

	result := <-getWalletListChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return 0, result.Error
	}

	return result.Wallet.CurrentSaldo, nil
}
