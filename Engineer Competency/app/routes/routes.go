package routes

import (
	"context"
	"engineer-comp/app/controllers"
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/middlewares"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func InitHTTPRoute(g *gin.Engine, databases pgsql.SqlInterface, ctx context.Context) {
	g.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	apiV1 := g.Group("/api/v1")
	{
		apiV1.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "initialized",
				"message": "API Gateway Internal Auth Service is running",
			})
		})
	}

	rootGroup := g.Group("api/v1", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTHBASIC_USERNAME"): os.Getenv("AUTHBASIC_PASSWORD"),
	}))
	oauthController := controllers.InitHTTPOauthController(databases)
	oauthGroup := rootGroup.Group("oauth")
	oauthGroup.Use()
	{
		oauthGroup.POST("token", oauthController.GenerateToken)
		oauthGroup.POST("verify-login-token", oauthController.VerifyAndValidateLoginToken)

	}
	//init controller
	walletController := controllers.InitHTTPWalletController(databases)
	invoiceController := controllers.InitHTTPInvoiceController(databases)

	oauthRootGroup := g.Group("api/v1")
	oauthRootGroup.Use(middlewares.OauthMiddleware(databases))
	{
		oauthRootGroup.GET("/protected-resource", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "You have accessed a protected resource",
			})
		})

		walletControllerGroup := oauthRootGroup.Group("/wallet")
		{
			walletControllerGroup.POST("/create", walletController.CreateWallet)
			walletControllerGroup.POST("/top-up", walletController.TopUpWallet)
			walletControllerGroup.GET("/:id", walletController.GetWalletByID)
			walletControllerGroup.GET("/user/:user_id", walletController.GetWalletByUserID)
			walletControllerGroup.GET("/balance/:user_id", walletController.GetWalletBalanceByUserID)
		}

		invoiceControllerGroup := oauthRootGroup.Group("/invoice")
		{
			invoiceControllerGroup.POST("/create", invoiceController.CreateInvoice)
			invoiceControllerGroup.GET("/:id", invoiceController.GetInvoiceByID)
			invoiceControllerGroup.GET("/all", invoiceController.GetAllInvoice)
			invoiceControllerGroup.PUT("/update-status", invoiceController.UpdateInvoice)

		}
	}

	_ = databases
	_ = ctx
}
