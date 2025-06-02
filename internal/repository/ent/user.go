package ent

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/omkar273/police/ent"
	"github.com/omkar273/police/ent/user"
	domainUser "github.com/omkar273/police/internal/domain/user"
	ierr "github.com/omkar273/police/internal/errors"
	"github.com/omkar273/police/internal/logger"
	"github.com/omkar273/police/internal/postgres"
	"github.com/omkar273/police/internal/types"
)

type UserRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts UserQueryOptions
}

func NewUserRepository(client postgres.IClient, logger *logger.Logger) domainUser.Repository {
	return &UserRepository{
		client:    client,
		log:       *logger,
		queryOpts: UserQueryOptions{},
	}
}

func (r *UserRepository) Create(ctx context.Context, userData *domainUser.User) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating user",
		"user_id", userData.ID,
		"email", userData.Email,
		"full_name", userData.FullName,
	)

	roles := make([]string, len(userData.Roles))
	for i, role := range userData.Roles {
		roles[i] = string(role)
	}

	_, err := client.User.Create().
		SetID(userData.ID).
		SetEmail(userData.Email).
		SetPhoneNumber(userData.Phone).
		SetFullName(userData.FullName).
		SetRoles(roles).
		SetStatus(string(userData.Status)).
		SetCreatedAt(userData.CreatedAt).
		SetUpdatedAt(userData.UpdatedAt).
		SetCreatedBy(userData.CreatedBy).
		SetUpdatedBy(userData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("User with this email already exists").
				WithReportableDetails(map[string]any{
					"user_id": userData.ID,
					"email":   userData.Email,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create user").
			WithReportableDetails(map[string]any{
				"user_id":   userData.ID,
				"user_name": userData.FullName,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (*domainUser.User, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting user", "user_id", id)

	entUser, err := client.User.Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("User with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"user_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get user").
			WithReportableDetails(map[string]any{
				"user_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainUser.FromEnt(entUser), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting user by email", "email", email)

	entUser, err := client.User.Query().
		Where(
			user.Email(email),
			user.StatusNotIn(string(types.StatusDeleted)),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("User with email %s was not found", email).
				WithReportableDetails(map[string]any{
					"email": email,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get user by email").
			WithReportableDetails(map[string]any{
				"email": email,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainUser.FromEnt(entUser), nil
}

func (r *UserRepository) List(ctx context.Context, filter *types.UserFilter) ([]*domainUser.User, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing users",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.User.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	users, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list users").
			Mark(ierr.ErrDatabase)
	}

	return domainUser.FromEntList(users), nil
}

func (r *UserRepository) ListAll(ctx context.Context, filter *types.UserFilter) ([]*domainUser.User, error) {
	if filter == nil {
		filter = types.NewNoLimitUserFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	users, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context, filter *types.UserFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting users")

	query := client.User.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count users").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *UserRepository) Update(ctx context.Context, userData *domainUser.User) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating user",
		"user_id", userData.ID,
		"email", userData.Email,
	)

	roles := make([]string, len(userData.Roles))
	for i, role := range userData.Roles {
		roles[i] = string(role)
	}

	_, err := client.User.UpdateOneID(userData.ID).
		SetEmail(userData.Email).
		SetPhoneNumber(userData.Phone).
		SetFullName(userData.FullName).
		SetRoles(roles).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(userData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("User with ID %s was not found", userData.ID).
				WithReportableDetails(map[string]any{
					"user_id": userData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("User with this email already exists").
				WithReportableDetails(map[string]any{
					"user_id": userData.ID,
					"email":   userData.Email,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to update user").
			WithReportableDetails(map[string]any{
				"user_id": userData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userData *domainUser.User) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting user",
		"user_id", userData.ID,
	)

	_, err := client.User.UpdateOneID(userData.ID).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(userData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("User with ID %s was not found", userData.ID).
				WithReportableDetails(map[string]any{
					"user_id": userData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete user").
			WithReportableDetails(map[string]any{
				"user_id": userData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

// UserQuery type alias for better readability
type UserQuery = *ent.UserQuery

// UserQueryOptions implements query options for user queries
type UserQueryOptions struct {
	QueryOptionsHelper
}

// Ensure UserQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[UserQuery, *types.UserFilter] = (*UserQueryOptions)(nil)

func (o UserQueryOptions) ApplyStatusFilter(query UserQuery, status string) UserQuery {
	if status == "" {
		return query.Where(user.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(user.Status(status))
}

func (o UserQueryOptions) ApplySortFilter(query UserQuery, field string, order string) UserQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o UserQueryOptions) ApplyPaginationFilter(query UserQuery, limit int, offset int) UserQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o UserQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return user.FieldCreatedAt
	case "updated_at":
		return user.FieldUpdatedAt
	case "email":
		return user.FieldEmail
	case "full_name":
		return user.FieldFullName
	case "phone":
		return user.FieldPhoneNumber
	default:
		return field
	}
}

func (o UserQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query UserQuery,
	filter *types.UserFilter,
) UserQuery {
	if filter == nil {
		return query.Where(user.StatusNotIn(string(types.StatusDeleted)))
	}

	// Apply status filter
	query = o.ApplyStatusFilter(query, filter.GetStatus())

	// Apply pagination
	if !filter.IsUnlimited() {
		query = o.ApplyPaginationFilter(query, filter.GetLimit(), filter.GetOffset())
	}

	// Apply sorting
	query = o.ApplySortFilter(query, filter.GetSort(), filter.GetOrder())

	return query
}

func (o UserQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.UserFilter,
	query UserQuery,
) UserQuery {
	if f == nil {
		return query
	}

	// Apply email filter if specified
	if len(f.Email) > 0 {
		query = query.Where(user.EmailIn(f.Email...))
	}

	// Apply phone filter if specified
	if len(f.Phone) > 0 {
		query = query.Where(user.PhoneNumberIn(f.Phone...))
	}

	// Apply roles filter if specified
	if len(f.Roles) > 0 {
		query = query.Where(func(s *sql.Selector) {
			predicates := make([]*sql.Predicate, len(f.Roles))
			for i, role := range f.Roles {
				predicates[i] = sql.And(
					sql.NotNull(s.C(user.FieldRoles)),
					sql.Contains(s.C(user.FieldRoles), role),
				)
			}
			s.Where(sql.And(predicates...))
		})

	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(user.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(user.CreatedAtLTE(*f.EndTime))
		}
	}

	return query
}
