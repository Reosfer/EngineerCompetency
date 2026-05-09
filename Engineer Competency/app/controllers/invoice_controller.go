package controllers

import (
	"engineer-comp/app/global/jwt"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/global/utils/model"
	"engineer-comp/app/models"
	usecases "engineer-comp/app/usecase"

	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	invoiceUseCase usecases.InvoiceUseCaseInterface
	oauthUsecase   usecases.OauthUseCaseInterface
}

func InitInvoiceController(
	invoiceUseCase usecases.InvoiceUseCaseInterface,
	oauthUsecase usecases.OauthUseCaseInterface,
) *InvoiceController {

	return &InvoiceController{
		invoiceUseCase: invoiceUseCase,
		oauthUsecase:   oauthUsecase,
	}
}

// CreateInvoice godoc
// @Summary Create invoice
// @Description Create a new invoice for a merchant
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body models.InvoiceRequest true "Invoice create request"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/invoice/create [post]
func (c *InvoiceController) CreateInvoice(ctx *gin.Context) {

	var result model.Response
	request := &models.InvoiceRequest{}
	err := ctx.BindJSON(request)
	if err != nil {
		fmt.Println("create invoice bind json error, " + err.Error())

		result.Error = err.Error()
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	if request.RoleID == "" {
		fmt.Println("role id is required")
		result.Error = "role id is required"
		result.StatusCode = http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	if request.RoleID != "2" {
		fmt.Println("merchant role id is required")
		result.Error = "merchant role id is required"
		result.StatusCode = http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	authorizationToken := ctx.Request.Header.Get("Authorization")

	if len(authorizationToken) == 0 {
		result.StatusCode = http.StatusUnauthorized
		result.Error = helper.NewError("Unauthorized").Error()
		ctx.Status(http.StatusUnauthorized)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, result)
		return
	}
	token := helper.GetAuthorizationValue(authorizationToken)
	tokenValue, err := jwt.TokenDecode(token)

	if err != nil {
		result.StatusCode = http.StatusUnauthorized
		result.Error = helper.NewError("Invalid token").Error()
		ctx.Status(http.StatusUnauthorized)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, result)
		return
	}
	//check user role
	_, _, err = c.oauthUsecase.CheckUserRole(tokenValue.Email, request.RoleID)
	if err != nil {
		fmt.Println("error checking user role, " + err.Error())
		result.StatusCode = http.StatusInternalServerError
		result.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	response, err := c.invoiceUseCase.CreateInvoice(request)

	if err != nil {

		fmt.Println("create invoice error, " + err.Error())

		if strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "greater than zero") {

			result.StatusCode = http.StatusBadRequest
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()

		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = response
	result.StatusCode = http.StatusCreated

	ctx.JSON(http.StatusCreated, result)
}

// GetInvoiceByID godoc
// @Summary Get invoice by ID
// @Description Get invoice details by invoice ID
// @Tags Invoice
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/invoice/{id} [get]
func (c *InvoiceController) GetInvoiceByID(ctx *gin.Context) {

	var result model.Response

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {

		fmt.Println("invalid invoice id")

		result.Error = "invalid invoice id"
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.invoiceUseCase.GetInvoiceByID(id)

	if err != nil {

		fmt.Println("get invoice by id error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") {
			result.StatusCode = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "not found") {
			result.StatusCode = http.StatusNotFound
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()

		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = response
	result.StatusCode = http.StatusOK

	ctx.JSON(http.StatusOK, result)
}

// GetAllInvoice godoc
// @Summary Get all invoices
// @Description Get list of all invoices
// @Tags Invoice
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/invoice/all [get]
func (c *InvoiceController) GetAllInvoice(ctx *gin.Context) {

	var result model.Response

	response, err := c.invoiceUseCase.GetAllInvoice()

	if err != nil {

		fmt.Println("get all invoice error, " + err.Error())

		result.Error = err.Error()
		result.StatusCode = http.StatusInternalServerError

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result.Data = response
	result.StatusCode = http.StatusOK

	ctx.JSON(http.StatusOK, result)
}

// UpdateInvoice godoc
// @Summary Update invoice status
// @Description Update invoice fields such as status and payment type
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body models.InvoiceRequest true "Invoice update request"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/invoice/update-status [put]
func (c *InvoiceController) UpdateInvoice(ctx *gin.Context) {

	var result model.Response

	request := &models.InvoiceRequest{}

	err := ctx.BindJSON(request)

	if err != nil {

		fmt.Println("update invoice bind json error, " + err.Error())

		result.Error = err.Error()
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	authorizationToken := ctx.Request.Header.Get("Authorization")

	if len(authorizationToken) == 0 {
		result.StatusCode = http.StatusUnauthorized
		result.Error = helper.NewError("Unauthorized").Error()
		ctx.Status(http.StatusUnauthorized)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, result)
		return
	}
	token := helper.GetAuthorizationValue(authorizationToken)
	tokenValue, err := jwt.TokenDecode(token)

	if err != nil {
		result.StatusCode = http.StatusUnauthorized
		result.Error = helper.NewError("Invalid token").Error()
		ctx.Status(http.StatusUnauthorized)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, result)
		return
	}
	//check user role
	_, userId, err := c.oauthUsecase.CheckUserRole(tokenValue.Email, request.RoleID)
	if err != nil {
		fmt.Println("error checking user role, " + err.Error())
		result.StatusCode = http.StatusInternalServerError
		result.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	println("user role id: " + request.RoleID)
	println("user id: " + fmt.Sprintf("%d", userId))

	if request.RoleID == "2" {
		request.MerchantID = userId
	} else if request.RoleID == "1" {
		request.BuyerId = userId
	} else if request.RoleID == "3" {
		request.AdminID = userId
	}

	_, err = c.invoiceUseCase.UpdateInvoice(request)

	if err != nil {
		fmt.Println("hah")
		fmt.Println("update invoice error, " + err.Error())

		if strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "greater than zero") {

			result.StatusCode = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "not found") {
			result.StatusCode = http.StatusNotFound
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = gin.H{
		"message": "invoice updated successfully",
	}
	result.StatusCode = http.StatusOK

	ctx.JSON(http.StatusOK, result)
}

func (c *InvoiceController) DeleteInvoice(ctx *gin.Context) {

	var result model.Response

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {

		fmt.Println("invalid invoice id")

		result.Error = "invalid invoice id"
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	err = c.invoiceUseCase.DeleteInvoice(id)

	if err != nil {

		fmt.Println("delete invoice error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") {
			result.StatusCode = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "not found") {
			result.StatusCode = http.StatusNotFound
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()

		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = "invoice deleted successfully"
	result.StatusCode = http.StatusOK

	ctx.JSON(http.StatusOK, result)
}
