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
	Health          *v1.HealthHandler
	Auth            *v1.AuthHandler
	User            *v1.UserHandler
	Internship      *v1.InternshipHandler
	InternshipBatch *v1.InternshipBatchHandler
	Category        *v1.CategoryHandler
	Discount        *v1.DiscountHandler
	Cart            *v1.CartHandler
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

	// Discount routes
	v1Discount := v1Router.Group("/discounts")
	{
		v1Discount.GET("", handlers.Discount.ListDiscounts)
		v1Discount.GET("/:id", handlers.Discount.GetDiscount)
		v1Discount.GET("/code/:code", handlers.Discount.GetDiscountByCode)

		v1Discount.Use(middleware.AuthenticateMiddleware(cfg, logger))
		v1Discount.POST("", handlers.Discount.CreateDiscount)
		v1Discount.PUT("/:id", handlers.Discount.UpdateDiscount)
		v1Discount.DELETE("/:id", handlers.Discount.DeleteDiscount)
	}

	// Cart routes
	v1Cart := v1Router.Group("/carts")
	v1Cart.Use(middleware.AuthenticateMiddleware(cfg, logger))
	{
		v1Cart.POST("", handlers.Cart.CreateCart)
		v1Cart.GET("/default", handlers.Cart.GetDefaultCart)
		v1Cart.GET("", handlers.Cart.ListCarts)
		v1Cart.GET("/:id", handlers.Cart.GetCart)
		v1Cart.PUT("/:id", handlers.Cart.UpdateCart)
		v1Cart.DELETE("/:id", handlers.Cart.DeleteCart)

		// Cart line items routes
		v1Cart.GET("/:id/line-items", handlers.Cart.GetCartLineItems)
		v1Cart.POST("/:id/line-items", handlers.Cart.AddLineItem)
		v1Cart.GET("/:id/line-items/:line_item_id", handlers.Cart.GetLineItem)
		v1Cart.DELETE("/:id/line-items/:line_item_id", handlers.Cart.RemoveLineItem)
	}

	// Internship Batch routes
	v1InternshipBatch := v1Router.Group("/internshipbatches")
	{
		v1InternshipBatch.GET("", handlers.InternshipBatch.ListInternshipBatches)
		v1InternshipBatch.GET("/:id", handlers.InternshipBatch.GetInternshipBatch)

		v1InternshipBatch.Use(middleware.AuthenticateMiddleware(cfg, logger))
		v1InternshipBatch.POST("", handlers.InternshipBatch.CreateInternshipBatch)
		v1InternshipBatch.PUT("/:id", handlers.InternshipBatch.UpdateInternshipBatch)
		v1InternshipBatch.DELETE("/:id", handlers.InternshipBatch.DeleteInternshipBatch)
	}

	return router
}
