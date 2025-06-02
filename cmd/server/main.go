package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/omkar273/codegeeky/docs/swagger"
	"github.com/omkar273/codegeeky/internal/api"
	v1 "github.com/omkar273/codegeeky/internal/api/v1"
	"github.com/omkar273/codegeeky/internal/auth"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/repository"
	"github.com/omkar273/codegeeky/internal/security"
	"github.com/omkar273/codegeeky/internal/service"
	"github.com/omkar273/codegeeky/internal/validator"
	"go.uber.org/fx"
)

// @title           CodeGeeky API
// @version         1.0
// @description     API for CodeGeeky
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. `Bearer <token>`.
// @type apiKey
// @required

func init() {
	// set time to UTC
	time.Local = time.UTC
}

func main() {
	var opts []fx.Option

	// load config
	opts = append(opts,
		fx.Provide(
			// provide config
			config.NewConfig,

			// validator
			validator.NewValidator,

			// logger
			logger.NewLogger,

			// postgres
			postgres.NewEntClient,
			postgres.NewClient,

			// auth provider
			auth.NewSupabaseProvider,

			// user repository
			repository.NewUserRepository,
		),
	)

	// services
	opts = append(opts, fx.Provide(

		// all services
		security.NewEncryptionService,
		service.NewAuthService,
		service.NewUserService,
		service.NewOnboardingService,
	))

	// factory layer
	opts = append(opts, fx.Provide(
		// handlers
		provideHandlers,

		// router
		provideRouter,
	))

	// start the application
	opts = append(opts, fx.Invoke(
		// start server
		startServer,
	))

	// start server
	app := fx.New(opts...)
	app.Run()
}

func startServer(
	lc fx.Lifecycle,
	cfg *config.Configuration,
	r *gin.Engine,
	log *logger.Logger,
) {
	// start api server
	startAPIServer(lc, r, cfg, log)
}

func provideHandlers(logger *logger.Logger, authService service.AuthService, userService service.UserService) *api.Handlers {
	return &api.Handlers{
		Health: v1.NewHealthHandler(logger),
		Auth:   v1.NewAuthHandler(authService),
		User:   v1.NewUserHandler(userService),
	}
}

func provideRouter(handlers *api.Handlers, cfg *config.Configuration, logger *logger.Logger) *gin.Engine {
	return api.NewRouter(handlers, cfg, logger)
}

func startAPIServer(
	lc fx.Lifecycle,
	r *gin.Engine,
	cfg *config.Configuration,
	log *logger.Logger,
) {
	log.Info("Registering API server start hook")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting API server...")
			go func() {
				if err := r.Run(cfg.Server.Address); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			log.Info("Server started successfully on port %s", cfg.Server.Address)
			log.Info("Server running at http://localhost%s", cfg.Server.Address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down server...")
			return nil
		},
	})
}
