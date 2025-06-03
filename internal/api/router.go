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
	Health     *v1.HealthHandler
	Auth       *v1.AuthHandler
	User       *v1.UserHandler
	Internship *v1.InternshipHandler
	Category   *v1.CategoryHandler
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
		v1Private.GET("/user/me", handlers.User.Me)
		v1Private.PUT("/user", handlers.User.Update)
	}

	// Internship routes
	v1Internship := v1Router.Group("/internships")
	{
		v1Internship.GET("", handlers.Internship.ListInternships)
		v1Internship.GET("/:id", handlers.Internship.GetInternship)

		v1Internship.Use(middleware.AuthenticateMiddleware(cfg, logger))
		v1Internship.POST("", handlers.Internship.CreateInternship)
		v1Internship.PUT("/:id", handlers.Internship.UpdateInternship)
		v1Internship.DELETE("/:id", handlers.Internship.DeleteInternship)
	}

	// Category routes
	v1Category := v1Router.Group("/categories")
	{
		v1Category.GET("", handlers.Category.ListCategories)
		v1Category.GET("/:id", handlers.Category.GetCategory)

		v1Category.Use(middleware.AuthenticateMiddleware(cfg, logger))
		v1Category.POST("", handlers.Category.CreateCategory)
		v1Category.PUT("/:id", handlers.Category.UpdateCategory)
		v1Category.DELETE("/:id", handlers.Category.DeleteCategory)
	}

	return router
}
