package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/omkar273/codegeeky/internal/api/v1"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/rest/middleware"
	"github.com/omkar273/codegeeky/internal/types"
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

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Global health check
	router.GET("/health", handlers.Health.Health)

	v1Router := router.Group("/v1")

	// Public routes
	v1Auth := v1Router.Group("/auth")
	v1Auth.Use(middleware.GuestAuthenticateMiddleware)
	v1Auth.POST("/signup", handlers.Auth.Signup)

	// Authenticated routes
	v1Private := v1Router.Group("/")
	v1Private.Use(middleware.AuthenticateMiddleware(cfg, logger))
	{
		v1Private.GET("/me", func(c *gin.Context) {
			userID := types.GetUserID(c.Request.Context())
			userEmail := types.GetUserEmail(c.Request.Context())
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "user_email": userEmail})
		})
	}

	return router
}
