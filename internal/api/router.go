package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/omkar273/codegeeky/internal/api/v1"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/rest/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Health *v1.HealthHandler
	Auth   *v1.AuthHandler
}

func NewRouter(handlers *Handlers, cfg *config.Configuration, logger *logger.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(
		middleware.CORSMiddleware,
		middleware.RequestIDMiddleware,
		middleware.ErrorHandler(),
	)

	// swagger// Add Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", handlers.Health.Health)
	router.POST("/health", handlers.Health.Health)

	v1Router := router.Group("/v1")
	v1Router.Use(middleware.ErrorHandler())
	{
		v1Router.GET("/health", handlers.Health.Health)
	}

	v1Auth := v1Router.Group("/auth")
	{
		v1Auth.POST("/signup", handlers.Auth.Signup)
	}

	return router
}
