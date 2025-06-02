package service

import (
	"github.com/omkar273/police/internal/config"
	"github.com/omkar273/police/internal/logger"
	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	// Core dependencies
	Logger *logger.Logger
	Config *config.Configuration

	// Repository dependencies
}
