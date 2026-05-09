package controllers

import (
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/global/utils/model"
	"engineer-comp/app/models"
	usecases "engineer-comp/app/usecase"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type OauthController struct {
	oauthUseCase usecases.OauthUseCaseInterface
}

func InitOauthController(oauthUseCase usecases.OauthUseCaseInterface) *OauthController {
	return &OauthController{
		oauthUseCase: oauthUseCase,
	}
}

// GenerateToken godoc
// @Summary Generate OAuth token
// @Description Generate OAuth access token using client credentials and password
// @Tags OAuth
// @Accept json
// @Produce json
// @Param request body models.Request true "OAuth token request"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BasicAuth
// @Router /api/v1/oauth/token [post]
func (c *OauthController) GenerateToken(ctx *gin.Context) {
	var result model.Response
	tokenRequest := &models.Request{}
	err := ctx.BindJSON(tokenRequest)

	if err != nil {
		fmt.Println("error token request, " + err.Error())
		result.Error = err.Error()
		result.StatusCode = http.StatusInternalServerError
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	var validationErrors []error

	if len(tokenRequest.ClientID) == 0 {
		fmt.Println("Client ID is Required")
		err = helper.NewError("Client ID is Required")
		validationErrors = append(validationErrors, err)
	}

	if len(tokenRequest.ClientSecret) == 0 {
		fmt.Println("Client Secret is Required")
		err = helper.NewError("Client Secret is Required")
		validationErrors = append(validationErrors, err)
	}
	if len(tokenRequest.Password) == 0 {
		fmt.Println("Password is Required")
		err = helper.NewError("Password is Required")
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		fmt.Println("error field request token")
		errorMessages := []string{}
		for _, v := range validationErrors {
			errorMessages = append(errorMessages, v.Error())
		}

		result.Error = errorMessages
		result.StatusCode = http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	tokenResponse, err := c.oauthUseCase.GenerateToken(tokenRequest)

	if err != nil {
		fmt.Println("GenerateToken error, " + err.Error())
		if strings.Contains(err.Error(), "token expired") {
			result.StatusCode = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			result.StatusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "Wrong password") {
			result.StatusCode = http.StatusUnauthorized
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = tokenResponse
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

// VerifyAndValidateLoginToken godoc
// @Summary Verify login token
// @Description Verify and validate an existing login token
// @Tags OAuth
// @Accept json
// @Produce json
// @Param request body models.Request true "Login token request"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Security BasicAuth
// @Router /api/v1/oauth/verify-login-token [post]
func (c *OauthController) VerifyAndValidateLoginToken(ctx *gin.Context) {
	var result model.Response
	tokenRequest := &models.Request{}
	err := ctx.BindJSON(tokenRequest)

	if err != nil {
		fmt.Println("tokenRequest bindjson error, " + err.Error())
		result.Error = err
		result.StatusCode = http.StatusInternalServerError
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	var validationErrors []error
	if len(tokenRequest.UserToken) == 0 {
		err = helper.NewError("Login Token is Required")
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		fmt.Println("Login Token error, required")
		errorMessages := []string{}
		for _, v := range validationErrors {
			errorMessages = append(errorMessages, v.Error())
		}

		result.Error = errorMessages
		result.StatusCode = http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validateTokenWithClient, err := c.oauthUseCase.VerifyAndValidateLoginToken(tokenRequest.UserToken)

	if err != nil {
		fmt.Println("validate token fail, error," + err.Error())
		if strings.Contains(err.Error(), "token expired") {
			result.StatusCode = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			result.StatusCode = http.StatusNotFound
		} else {
			result.StatusCode = http.StatusInternalServerError
		}

		result.Error = err.Error()
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = validateTokenWithClient
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}
