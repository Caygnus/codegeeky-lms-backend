package v1

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/omkar273/codegeeky/internal/auth"
// 	"github.com/omkar273/codegeeky/internal/config"
// 	domainAuth "github.com/omkar273/codegeeky/internal/domain/auth"
// 	domainUser "github.com/omkar273/codegeeky/internal/domain/user"
// 	"github.com/omkar273/codegeeky/internal/logger"
// 	"github.com/omkar273/codegeeky/internal/rest/middleware"
// 	"github.com/omkar273/codegeeky/internal/service"
// )

// // SetupInternshipRoutes sets up the internship routes with proper authorization
// func SetupInternshipRoutes(
// 	router *gin.RouterGroup,
// 	cfg *config.Configuration,
// 	logger *logger.Logger,
// 	authzService auth.AuthorizationService,
// 	userRepo domainUser.Repository,
// 	internshipService service.InternshipService,
// ) {
// 	// Create handler
// 	handler := NewInternshipHandler(internshipService, authzService, logger)

// 	// Set up middleware for authorization
// 	authMiddleware := middleware.AuthorizationMiddleware(cfg, logger, authzService, userRepo)

// 	// Apply authentication and authorization middleware to all routes
// 	internshipGroup := router.Group("/internships")
// 	internshipGroup.Use(authMiddleware)

// 	// Public routes (all authenticated users can access)
// 	internshipGroup.GET("", handler.ListInternships) // Role-based filtering applied in service

// 	// RBAC: Only instructors and admins can create internships
// 	internshipGroup.POST("",
// 		middleware.RequireInstructorOrAdmin(),
// 		handler.CreateInternship,
// 	)

// 	// ABAC: Fine-grained access control for viewing specific internships
// 	internshipGroup.GET("/:id", handler.GetInternship)

// 	// ABAC: Ownership-based access control for updating internships
// 	internshipGroup.PUT("/:id", handler.UpdateInternship)
// 	internshipGroup.PATCH("/:id", handler.UpdateInternship)

// 	// ABAC: Ownership-based access control for deleting internships
// 	internshipGroup.DELETE("/:id", handler.DeleteInternship)

// 	// ABAC: Content access with enrollment and progress checks
// 	internshipGroup.GET("/:internship_id/content/:content_id", handler.AccessInternshipContent)
// }

// // SetupRoutesWithPermissionMiddleware demonstrates using permission-based middleware
// func SetupRoutesWithPermissionMiddleware(
// 	router *gin.RouterGroup,
// 	cfg *config.Configuration,
// 	logger *logger.Logger,
// 	authzService auth.AuthorizationService,
// 	userRepo domainUser.Repository,
// 	internshipService service.InternshipService,
// ) {
// 	// Create handler
// 	handler := NewInternshipHandler(internshipService, authzService, logger)

// 	// Set up middleware for authorization
// 	authMiddleware := middleware.AuthorizationMiddleware(cfg, logger, authzService, userRepo)

// 	// Apply authentication and authorization middleware to all routes
// 	internshipGroup := router.Group("/internships")
// 	internshipGroup.Use(authMiddleware)

// 	// Using permission-based middleware instead of role-based
// 	internshipGroup.POST("",
// 		middleware.RequirePermission(domainAuth.PermissionCreateInternship, "internship"),
// 		handler.CreateInternship,
// 	)

// 	internshipGroup.GET("/:id",
// 		middleware.RequirePermission(domainAuth.PermissionViewInternship, "internship"),
// 		handler.GetInternship,
// 	)

// 	internshipGroup.PUT("/:id",
// 		middleware.RequirePermission(domainAuth.PermissionUpdateInternship, "internship"),
// 		handler.UpdateInternship,
// 	)

// 	internshipGroup.DELETE("/:id",
// 		middleware.RequirePermission(domainAuth.PermissionDeleteInternship, "internship"),
// 		handler.DeleteInternship,
// 	)

// 	// Content access routes with specific permissions
// 	contentGroup := internshipGroup.Group("/:internship_id/content")

// 	contentGroup.GET("/lectures/:content_id",
// 		middleware.RequirePermission(domainAuth.PermissionViewLectures, "content"),
// 		handler.AccessInternshipContent,
// 	)

// 	contentGroup.GET("/assignments/:content_id",
// 		middleware.RequirePermission(domainAuth.PermissionViewAssignments, "content"),
// 		handler.AccessInternshipContent,
// 	)

// 	contentGroup.GET("/resources/:content_id",
// 		middleware.RequirePermission(domainAuth.PermissionViewResources, "content"),
// 		handler.AccessInternshipContent,
// 	)
// }

// // SetupAdminRoutes demonstrates admin-only routes
// func SetupAdminRoutes(
// 	router *gin.RouterGroup,
// 	cfg *config.Configuration,
// 	logger *logger.Logger,
// 	authzService auth.AuthorizationService,
// 	userRepo domainUser.Repository,
// 	internshipService service.InternshipService,
// ) {
// 	// Create handler
// 	handler := NewInternshipHandler(internshipService, authzService, logger)

// 	// Set up middleware for authorization
// 	authMiddleware := middleware.AuthorizationMiddleware(cfg, logger, authzService, userRepo)

// 	// Admin-only routes
// 	adminGroup := router.Group("/admin")
// 	adminGroup.Use(authMiddleware)
// 	adminGroup.Use(middleware.RequireAdmin())

// 	// Admin can see all internships regardless of status
// 	adminGroup.GET("/internships", handler.ListInternships)

// 	// Admin can manage system configuration
// 	adminGroup.GET("/analytics",
// 		middleware.RequirePermission(domainAuth.PermissionViewAnalytics, "system"),
// 		func(c *gin.Context) {
// 			c.JSON(200, gin.H{"message": "Analytics data"})
// 		},
// 	)

// 	adminGroup.POST("/system/config",
// 		middleware.RequirePermission(domainAuth.PermissionSystemConfig, "system"),
// 		func(c *gin.Context) {
// 			c.JSON(200, gin.H{"message": "System configuration updated"})
// 		},
// 	)
// }

// // Example of how to set up routes in your main router
// func ExampleRouterSetup(
// 	cfg *config.Configuration,
// 	logger *logger.Logger,
// 	authzService auth.AuthorizationService,
// 	userRepo domainUser.Repository,
// 	internshipService service.InternshipService,
// ) *gin.Engine {
// 	router := gin.Default()

// 	// Add global middleware
// 	router.Use(middleware.AuthenticateMiddleware(cfg, logger))

// 	// API v1 routes
// 	v1 := router.Group("/api/v1")

// 	// Set up internship routes
// 	SetupInternshipRoutes(v1, cfg, logger, authzService, userRepo, internshipService)

// 	// Set up admin routes
// 	SetupAdminRoutes(v1, cfg, logger, authzService, userRepo, internshipService)

// 	return router
// }
