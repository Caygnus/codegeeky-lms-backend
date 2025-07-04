package repository

import (
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/discount"
	"github.com/omkar273/codegeeky/internal/domain/internship"
	"github.com/omkar273/codegeeky/internal/domain/internshipenrollment"
	"github.com/omkar273/codegeeky/internal/domain/payment"
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

func NewInternshipRepository(params RepositoryParams) internship.InternshipRepository {
	return ent.NewInternshipRepository(params.Client, params.Logger)
}

func NewCategoryRepository(params RepositoryParams) internship.CategoryRepository {
	return ent.NewCategoryRepository(params.Client, params.Logger)
}

func NewDiscountRepository(params RepositoryParams) discount.Repository {
	return ent.NewDiscountRepository(params.Client, params.Logger)
}

func NewPaymentRepository(params RepositoryParams) payment.Repository {
	return ent.NewPaymentRepository(params.Client, params.Logger)
}

func NewInternshipEnrollmentRepository(params RepositoryParams) internshipenrollment.Repository {
	return ent.NewInternshipEnrollmentRepository(params.Client, params.Logger)
}

func NewInternshipBatchRepository(params RepositoryParams) internship.InternshipBatchRepository {
	return ent.NewInternshipBatchRepository(params.Client, params.Logger)
}
