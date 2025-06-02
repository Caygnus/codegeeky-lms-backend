package repository

import (
	"github.com/omkar273/police/internal/config"
	"github.com/omkar273/police/internal/logger"
	"github.com/omkar273/police/internal/postgres"
	"go.uber.org/fx"
)

type RepositoryParams struct {
	// factory params
	fx.In

	Client postgres.IClient
	Logger *logger.Logger
	Config *config.Configuration
}

func NewRepositoryParams(
	Client *postgres.Client,
	Logger *logger.Logger,
	Config *config.Configuration,
) RepositoryParams {
	return RepositoryParams{
		Client: Client,
		Logger: Logger,
	}
}
