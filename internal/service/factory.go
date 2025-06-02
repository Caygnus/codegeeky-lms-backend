package service

import (
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/user"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	// Core dependencies
	Logger *logger.Logger
	Config *config.Configuration
	DB     postgres.IClient

	// Repository dependencies
	UserRepo user.Repository
}
