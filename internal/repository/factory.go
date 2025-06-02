package repository

import (
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/user"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/repository/ent"
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

func NewUserRepository(params RepositoryParams) user.Repository {
	return ent.NewUserRepository(params.Client, params.Logger)
}
