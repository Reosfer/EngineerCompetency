package middlewares

import (
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"engineer-comp/app/repositories"
	usecases "engineer-comp/app/usecase"

	"net/http"

	"github.com/gin-gonic/gin"
)

func OauthMiddleware(databases pgsql.SqlInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result models.Response
		grantAccessRepository := repositories.InitGrantRepository(databases)
		useCase := usecases.InitOauthUseCase(grantAccessRepository)
		authorizationToken := c.Request.Header.Get("Authorization")

		if len(authorizationToken) == 0 {
			result.StatusCode = http.StatusUnauthorized
			result.Error = helper.NewError("Unauthorized").Error()
			c.Status(http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, result)
			return
		}

		token := helper.GetAuthorizationValue(authorizationToken)
		requestID, _ := c.Get("RequestId")

		resultOauthValidateResponseChan := make(chan *models.ValidUserTokenWithClient)
		go useCase.ValidateTokenUser(token, requestID.(string), resultOauthValidateResponseChan)
		resultOauthValidateResponse := <-resultOauthValidateResponseChan

		if resultOauthValidateResponse.Error != nil {
			result.StatusCode = resultOauthValidateResponse.StatusCode
			result.Error = resultOauthValidateResponse.Error.Error()
			c.Status(resultOauthValidateResponse.StatusCode)
			c.AbortWithStatusJSON(resultOauthValidateResponse.StatusCode, result)
			return
		}
		c.Set("LoginToken", token)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
