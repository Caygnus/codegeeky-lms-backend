package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/cache"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/cart"
	"github.com/omkar273/codegeeky/internal/domain/discount"
	"github.com/omkar273/codegeeky/internal/domain/internship"
	"github.com/omkar273/codegeeky/internal/domain/internshipenrollment"
	"github.com/omkar273/codegeeky/internal/domain/user"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/stretchr/testify/suite"
)

// Stores holds all the repository interfaces for testing
type Stores struct {
	UserRepo                 user.Repository
	CartRepo                 cart.Repository
	DiscountRepo             discount.Repository
	InternshipRepo           internship.InternshipRepository
	InternshipBatchRepo      internship.InternshipBatchRepository
	InternshipEnrollmentRepo internshipenrollment.Repository
}

// BaseServiceTestSuite provides common functionality for all service test suites
type BaseServiceTestSuite struct {
	suite.Suite
	ctx    context.Context
	stores Stores
	db     postgres.IClient
	logger *logger.Logger
	config *config.Configuration
	now    time.Time
}

// SetupSuite is called once before running the tests in the suite
func (s *BaseServiceTestSuite) SetupSuite() {
	// Initialize validator
	validator.NewValidator()

	// Initialize logger with test config
	cfg := &config.Configuration{
		Logging: config.LoggingConfig{
			Level: types.LogLevelDebug,
		},
	}
	var err error
	s.config = cfg
	s.logger, err = logger.NewLogger(cfg)
	if err != nil {
		s.T().Fatalf("failed to create logger: %v", err)
	}

	// Initialize cache
	cache.Initialize(s.logger)
}

// SetupTest is called before each test
func (s *BaseServiceTestSuite) SetupTest() {
	s.setupContext()
	s.setupStores()
	s.now = time.Now().UTC()
}

// TearDownTest is called after each test
func (s *BaseServiceTestSuite) TearDownTest() {
	s.clearStores()
}

func (s *BaseServiceTestSuite) setupContext() {
	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, types.CtxUserID, types.DefaultUserID)
	s.ctx = context.WithValue(s.ctx, types.CtxRequestID, types.GenerateUUID())
	s.ctx = context.WithValue(s.ctx, types.CtxUserEmail, types.DefaultUserEmail)
	s.ctx = context.WithValue(s.ctx, types.CtxUserRole, types.DefaultUserRole)
	s.ctx = context.WithValue(s.ctx, types.CtxUser, types.DefaultUserID)
	s.ctx = context.WithValue(s.ctx, types.CtxIsGuest, types.DefaultIsGuest)
}

func (s *BaseServiceTestSuite) setupStores() {
	s.stores = Stores{
		UserRepo:                 NewInMemoryUserStore(),
		CartRepo:                 NewInMemoryCartStore(),
		DiscountRepo:             NewInMemoryDiscountStore(),
		InternshipRepo:           NewInMemoryInternshipStore(),
		InternshipBatchRepo:      NewInMemoryInternshipBatchStore(),
		InternshipEnrollmentRepo: NewInMemoryInternshipEnrollmentStore(),
	}

	s.db = NewMockPostgresClient(s.logger)
}

func (s *BaseServiceTestSuite) clearStores() {
	s.stores.CartRepo.(*InMemoryCartStore).Clear()
	s.stores.UserRepo.(*InMemoryUserStore).Clear()
	s.stores.DiscountRepo.(*InMemoryDiscountStore).Clear()
	s.stores.InternshipRepo.(*InMemoryInternshipStore).Clear()
	s.stores.InternshipBatchRepo.(*InMemoryInternshipBatchStore).Clear()
	s.stores.InternshipEnrollmentRepo.(*InMemoryInternshipEnrollmentStore).Clear()
}

func (s *BaseServiceTestSuite) ClearStores() {
	s.clearStores()
}

// GetContext returns the test context
func (s *BaseServiceTestSuite) GetContext() context.Context {
	return s.ctx
}

// GetConfig returns the test configuration
func (s *BaseServiceTestSuite) GetConfig() *config.Configuration {
	return s.config
}

// GetStores returns all test repositories
func (s *BaseServiceTestSuite) GetStores() Stores {
	return s.stores
}

// GetDB returns the test database client
func (s *BaseServiceTestSuite) GetDB() postgres.IClient {
	return s.db
}

// GetLogger returns the test logger
func (s *BaseServiceTestSuite) GetLogger() *logger.Logger {
	return s.logger
}

// GetNow returns the current test time
func (s *BaseServiceTestSuite) GetNow() time.Time {
	return s.now.UTC()
}

// GetUUID returns a new UUID string
func (s *BaseServiceTestSuite) GetUUID() string {
	return types.GenerateUUID()
}
