package controllers

import (
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/repositories"
	usecases "engineer-comp/app/usecase"
)

func InitHTTPOauthController(databases pgsql.SqlInterface) *OauthController {
	grantAccessRepository := repositories.InitGrantRepository(databases)
	oauthUseCase := usecases.InitOauthUseCase(grantAccessRepository)
	handler := InitOauthController(oauthUseCase)
	return handler
}

func InitHTTPWalletController(databases pgsql.SqlInterface) *WalletController {
	walletRepository := repositories.InitWalletRepository(databases)
	walletUseCase := usecases.InitWalletUseCase(walletRepository)
	handler := InitWalletController(walletUseCase)
	return handler
}

func InitHTTPInvoiceController(databases pgsql.SqlInterface) *InvoiceController {
	invoiceRepository := repositories.InitInvoiceRepository(databases)
	invoiceUseCase := usecases.InitInvoiceUseCase(invoiceRepository)
	grantAccessRepository := repositories.InitGrantRepository(databases)
	oauthUseCase := usecases.InitOauthUseCase(grantAccessRepository)
	handler := InitInvoiceController(invoiceUseCase, oauthUseCase)
	return handler
}
