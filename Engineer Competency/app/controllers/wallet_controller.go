package controllers

import (
	"engineer-comp/app/global/utils/model"
	"engineer-comp/app/models"
	usecases "engineer-comp/app/usecase"

	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletUseCase usecases.WalletUseCaseInterface
}

func InitWalletController(
	walletUseCase usecases.WalletUseCaseInterface,
) *WalletController {

	return &WalletController{
		walletUseCase: walletUseCase,
	}
}

// CreateWallet godoc
// @Summary Create wallet
// @Description Create a new wallet for a user
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body models.Wallet true "Wallet request"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/wallet/create [post]
func (c *WalletController) CreateWallet(ctx *gin.Context) {

	var result model.Response

	request := &models.Wallet{}

	err := ctx.BindJSON(request)

	if err != nil {

		fmt.Println("create wallet bind json error, " + err.Error())

		result.Error = err.Error()
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.walletUseCase.CreateWallet(request)

	if err != nil {

		fmt.Println("create wallet error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "negative") {

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

// TopUpWallet godoc
// @Summary Top up wallet balance
// @Description Add balance to a wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body models.Wallet true "Top up request"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/wallet/top-up [post]
func (c *WalletController) TopUpWallet(ctx *gin.Context) {

	var result model.Response

	request := &models.Wallet{}
	err := ctx.BindJSON(request)

	if err != nil {

		fmt.Println("create wallet bind json error, " + err.Error())

		result.Error = err.Error()
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.walletUseCase.TopUpWallet(request)

	if err != nil {

		fmt.Println("create wallet error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "negative") {

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

// GetWalletByID godoc
// @Summary Get wallet by ID
// @Description Get wallet details by wallet ID
// @Tags Wallet
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/wallet/{id} [get]
func (c *WalletController) GetWalletByID(ctx *gin.Context) {

	var result model.Response

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {

		fmt.Println("invalid wallet id")

		result.Error = "invalid wallet id"
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.walletUseCase.GetWalletByID(id)

	if err != nil {

		fmt.Println("get wallet by id error, " + err.Error())

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

// GetWalletByUserID godoc
// @Summary Get wallet by user ID
// @Description Get wallet list by user ID
// @Tags Wallet
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/wallet/user/{user_id} [get]
func (c *WalletController) GetWalletByUserID(ctx *gin.Context) {

	var result model.Response

	userIDParam := ctx.Param("user_id")

	userID, err := strconv.Atoi(userIDParam)

	if err != nil {

		fmt.Println("invalid user id")

		result.Error = "invalid user id"
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.walletUseCase.GetWalletByUserID(userID)

	if err != nil {

		fmt.Println("get wallet by user id error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") {
			result.StatusCode = http.StatusBadRequest
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

// GetWalletBalanceByUserID godoc
// @Summary Get wallet balance by user ID
// @Description Get balance for a user wallet by user ID
// @Tags Wallet
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BearerAuth
// @Router /api/v1/wallet/balance/{user_id} [get]
func (c *WalletController) GetWalletBalanceByUserID(ctx *gin.Context) {

	var result model.Response

	userIDParam := ctx.Param("user_id")

	userID, err := strconv.Atoi(userIDParam)

	if err != nil {

		fmt.Println("invalid user id")

		result.Error = "invalid user id"
		result.StatusCode = http.StatusBadRequest

		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	response, err := c.walletUseCase.GetBalanceWalletByUserID(userID)

	if err != nil {

		fmt.Println("get wallet by user id error, " + err.Error())

		if strings.Contains(err.Error(), "invalid") {
			result.StatusCode = http.StatusBadRequest
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()

		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = gin.H{
		"user_id": userID,
		"balance": response,
	}
	result.StatusCode = http.StatusOK

	ctx.JSON(http.StatusOK, result)
}
