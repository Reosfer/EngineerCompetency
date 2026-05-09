// @title Engineer Competency API
// @version 1.0
// @description API Gateway Internal Auth Service dengan OAuth, wallet, dan invoice.
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	docs "engineer-comp/app/docs"
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/middlewares"
	"engineer-comp/app/routes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	envConfig "github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load .env file if it exists (for local development).
	// If .env is missing, try .env.docker.
	if err := envConfig.Load(".env"); err != nil {
		fmt.Println(".env is not loaded, trying .env.docker")
		if err2 := envConfig.Load(".env.docker"); err2 != nil {
			fmt.Println(".env.docker is not loaded")
			fmt.Println(err2)
		}
	}

	fmt.Println(os.Getenv("AUTHBASIC_USERNAME"), os.Getenv("AUTHBASIC_PASSWORD"), "hehe")

	fmt.Println("Starting API Gateway Internal Auth Service HTTP Handler")

	database, err := pgsql.InitPgSql(
		os.Getenv("DB_PGSQL_PROXY_URL"),
		"postgres",
		os.Getenv("DB_PGSQL_PROXY_HOST"),
		os.Getenv("DB_PGSQL_PROXY_PORT"),
		os.Getenv("DB_PGSQL_PROXY_DWH_USERNAME"),
		os.Getenv("DB_PGSQL_PROXY_DWH_PASSWORD"),
		os.Getenv("DB_PGSQL_PROXY_DWH_DATABASE"),
	)
	if err != nil {
		panic(err)
	}
	defer database.DB().Close()

	g := gin.Default()
	g.Use(middlewares.CORSMiddleware(), middlewares.JSONMiddleware(), RequestId())
	routes.InitHTTPRoute(g, database, context.Background())

	port := os.Getenv("MAIN_PORT")
	if port == "" {
		port = "8080"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Engineer Competency API"
	docs.SwaggerInfo.Description = "API Gateway Internal Auth Service dengan OAuth, wallet, dan invoice."

	swaggerURL := "http://localhost:" + port + "/swagger/doc.json"
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerURL)))

	addr := ":" + port
	if addr == "" {
		addr = ":8080"
	} else if addr[0] != ':' {
		addr = ":" + addr
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      g,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("RequestId", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
