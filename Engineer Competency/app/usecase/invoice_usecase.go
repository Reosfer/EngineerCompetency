package usecases

import (
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"engineer-comp/app/repositories"
	"strings"

	"fmt"
	"net/http"
)

type InvoiceUseCaseInterface interface {
	CreateInvoice(request *models.InvoiceRequest) (*models.Invoice, error)
	GetInvoiceByID(id int) (*models.Invoice, error)
	GetAllInvoice() ([]models.Invoice, error)
	UpdateInvoice(request *models.InvoiceRequest) (*models.Invoice, error)
	DeleteInvoice(id int) error
}

type invoiceUseCase struct {
	invoiceRepository repositories.InvoiceRepositoryInterface
}

func InitInvoiceUseCase(
	invoiceRepository repositories.InvoiceRepositoryInterface,
) InvoiceUseCaseInterface {

	return &invoiceUseCase{
		invoiceRepository: invoiceRepository,
	}
}

func (u *invoiceUseCase) CreateInvoice(request *models.InvoiceRequest) (*models.Invoice, error) {

	if request.InvoiceName == "" {
		return nil, helper.NewError("invoice name is required")
	}

	if request.Amount <= 0 {
		return nil, helper.NewError("amount must be greater than zero")
	}

	if request.Status == "" {
		request.Status = "pending"
	}

	validStatus := map[string]bool{
		"pending":  true,
		"approved": true,
		"refund":   true,
	}

	if !validStatus[strings.ToLower(request.Status)] {
		return nil, helper.NewError("invalid invoice status")
	}

	createInvoiceChan := make(chan *models.InvoiceResponseChan)

	go u.invoiceRepository.CreateInvoice(request, createInvoiceChan)

	result := <-createInvoiceChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Invoice, nil
}

func (u *invoiceUseCase) GetInvoiceByID(id int) (*models.Invoice, error) {

	if id <= 0 {
		return nil, helper.NewError("invalid invoice id")
	}

	getInvoiceChan := make(chan *models.InvoiceResponseChan)

	go u.invoiceRepository.GetInvoiceByID(id, getInvoiceChan)

	result := <-getInvoiceChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Invoice, nil
}

func (u *invoiceUseCase) GetAllInvoice() ([]models.Invoice, error) {

	getAllInvoiceChan := make(chan *models.InvoiceListResponseChan)

	go u.invoiceRepository.GetAllInvoice(getAllInvoiceChan)

	result := <-getAllInvoiceChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, result.Error
	}

	return result.Invoice, nil
}

func (u *invoiceUseCase) UpdateInvoice(request *models.InvoiceRequest) (*models.Invoice, error) {

	if request.ID <= 0 {
		return nil, helper.NewError("invalid invoice id")
	}

	validStatus := map[string]bool{
		"pending":  true,
		"approved": true,
		"refund":   true,
		"rejected": true,
		"paid":     true,
	}

	if !validStatus[request.Status] {
		return nil, helper.NewError("invalid invoice status")
	}

	updateInvoiceChan := make(chan *models.InvoiceResponseChan)

	go u.invoiceRepository.UpdateInvoice(request, updateInvoiceChan)

	result := <-updateInvoiceChan
	if result.Error != nil {
		fmt.Println(result.Error.Error(), "err")
		return nil, result.Error
	}

	return result.Invoice, nil
}

func (u *invoiceUseCase) DeleteInvoice(id int) error {

	if id <= 0 {
		return helper.NewError("invalid invoice id")
	}

	deleteInvoiceChan := make(chan *models.InvoiceResponseChan)

	go u.invoiceRepository.DeleteInvoice(id, deleteInvoiceChan)

	result := <-deleteInvoiceChan

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return result.Error
	}

	return nil
}

func InvoiceSuccessResponse(data interface{}) *models.Response {

	return &models.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
	}
}
