package ent

// import (
// 	"context"
// 	"time"

// 	"github.com/omkar273/codegeeky/ent"
// 	"github.com/omkar273/codegeeky/ent/rank"
// 	domainRank "github.com/omkar273/codegeeky/internal/domain/rank"
// 	ierr "github.com/omkar273/codegeeky/internal/errors"
// 	"github.com/omkar273/codegeeky/internal/logger"
// 	"github.com/omkar273/codegeeky/internal/postgres"
// 	"github.com/omkar273/codegeeky/internal/types"
// )

// type RankRepository struct {
// 	client    postgres.IClient
// 	log       logger.Logger
// 	queryOpts RankQueryOptions
// }

// func NewRankRepository(client postgres.IClient, logger *logger.Logger) domainRank.Repository {
// 	return &RankRepository{
// 		client:    client,
// 		log:       *logger,
// 		queryOpts: RankQueryOptions{},
// 	}
// }

// func (r *RankRepository) Create(ctx context.Context, rankData *domainRank.Rank) (*domainRank.Rank, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("creating rank",
// 		"rank_id", rankData.ID,
// 		"name", rankData.Name,
// 		"code", rankData.Code,
// 	)

// 	entRank, err := client.Rank.Create().
// 		SetName(rankData.Name).
// 		SetDescription(rankData.Description).
// 		SetAbbreviation(rankData.Abbreviation).
// 		SetCode(rankData.Code).
// 		SetHierarchyLevel(rankData.HierarchyLevel).
// 		SetStatus(string(rankData.Status)).
// 		SetCreatedAt(rankData.CreatedAt).
// 		SetUpdatedAt(rankData.UpdatedAt).
// 		SetCreatedBy(rankData.CreatedBy).
// 		SetUpdatedBy(rankData.UpdatedBy).
// 		Save(ctx)

// 	if err != nil {
// 		if ent.IsConstraintError(err) {
// 			return nil, ierr.WithError(err).
// 				WithHint("Rank with this code already exists").
// 				WithReportableDetails(map[string]any{
// 					"rank_id":   rankData.ID,
// 					"rank_code": rankData.Code,
// 				}).
// 				Mark(ierr.ErrAlreadyExists)
// 		}
// 		return nil, ierr.WithError(err).
// 			WithHint("Failed to create rank").
// 			WithReportableDetails(map[string]any{
// 				"rank_id":   rankData.ID,
// 				"rank_name": rankData.Name,
// 			}).
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return domainRank.FromEnt(entRank), nil
// }

// func (r *RankRepository) Get(ctx context.Context, id string) (*domainRank.Rank, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("getting rank", "rank_id", id)

// 	entRank, err := client.Rank.Query().
// 		Where(rank.ID(id)).
// 		Only(ctx)

// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return nil, ierr.WithError(err).
// 				WithHintf("Rank with ID %s was not found", id).
// 				WithReportableDetails(map[string]any{
// 					"rank_id": id,
// 				}).
// 				Mark(ierr.ErrNotFound)
// 		}
// 		return nil, ierr.WithError(err).
// 			WithHint("Failed to get rank").
// 			WithReportableDetails(map[string]any{
// 				"rank_id": id,
// 			}).
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return domainRank.FromEnt(entRank), nil
// }

// func (r *RankRepository) GetByCode(ctx context.Context, code string) (*domainRank.Rank, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("getting rank by code", "code", code)

// 	entRank, err := client.Rank.Query().
// 		Where(
// 			rank.Code(code),
// 			rank.StatusNotIn(string(types.StatusDeleted)),
// 		).
// 		Only(ctx)

// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return nil, ierr.WithError(err).
// 				WithHintf("Rank with code %s was not found", code).
// 				WithReportableDetails(map[string]any{
// 					"rank_code": code,
// 				}).
// 				Mark(ierr.ErrNotFound)
// 		}
// 		return nil, ierr.WithError(err).
// 			WithHint("Failed to get rank by code").
// 			WithReportableDetails(map[string]any{
// 				"rank_code": code,
// 			}).
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return domainRank.FromEnt(entRank), nil
// }

// func (r *RankRepository) List(ctx context.Context, filter *types.RankFilter) ([]*domainRank.Rank, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("listing ranks",
// 		"limit", filter.GetLimit(),
// 		"offset", filter.GetOffset(),
// 	)

// 	query := client.Rank.Query()
// 	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
// 	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

// 	ranks, err := query.All(ctx)
// 	if err != nil {
// 		return nil, ierr.WithError(err).
// 			WithHint("Failed to list ranks").
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return domainRank.FromEntList(ranks), nil
// }

// func (r *RankRepository) ListAll(ctx context.Context, filter *types.RankFilter) ([]*domainRank.Rank, error) {
// 	if filter == nil {
// 		filter = types.NewNoLimitRankFilter()
// 	}

// 	if filter.QueryFilter == nil {
// 		filter.QueryFilter = types.NewNoLimitQueryFilter()
// 	}

// 	ranks, err := r.List(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ranks, nil
// }

// func (r *RankRepository) Count(ctx context.Context, filter *types.RankFilter) (int, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("counting ranks")

// 	query := client.Rank.Query()
// 	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
// 	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

// 	count, err := query.Count(ctx)
// 	if err != nil {
// 		return 0, ierr.WithError(err).
// 			WithHint("Failed to count ranks").
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return count, nil
// }

// func (r *RankRepository) Update(ctx context.Context, rankData *domainRank.Rank) (*domainRank.Rank, error) {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("updating rank",
// 		"rank_id", rankData.ID,
// 		"name", rankData.Name,
// 	)

// 	entRank, err := client.Rank.UpdateOneID(rankData.ID).
// 		SetName(rankData.Name).
// 		SetDescription(rankData.Description).
// 		SetAbbreviation(rankData.Abbreviation).
// 		SetCode(rankData.Code).
// 		SetHierarchyLevel(rankData.HierarchyLevel).
// 		SetUpdatedAt(time.Now().UTC()).
// 		SetUpdatedBy(rankData.UpdatedBy).
// 		Save(ctx)

// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return nil, ierr.WithError(err).
// 				WithHintf("Rank with ID %s was not found", rankData.ID).
// 				WithReportableDetails(map[string]any{
// 					"rank_id": rankData.ID,
// 				}).
// 				Mark(ierr.ErrNotFound)
// 		}
// 		if ent.IsConstraintError(err) {
// 			return nil, ierr.WithError(err).
// 				WithHint("Rank with this code already exists").
// 				WithReportableDetails(map[string]any{
// 					"rank_id":   rankData.ID,
// 					"rank_code": rankData.Code,
// 				}).
// 				Mark(ierr.ErrAlreadyExists)
// 		}
// 		return nil, ierr.WithError(err).
// 			WithHint("Failed to update rank").
// 			WithReportableDetails(map[string]any{
// 				"rank_id": rankData.ID,
// 			}).
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return domainRank.FromEnt(entRank), nil
// }

// func (r *RankRepository) Delete(ctx context.Context, rankData *domainRank.Rank) error {
// 	client := r.client.Querier(ctx)

// 	r.log.Debugw("deleting rank",
// 		"rank_id", rankData.ID,
// 	)

// 	_, err := client.Rank.UpdateOneID(rankData.ID).
// 		SetStatus(string(types.StatusDeleted)).
// 		SetUpdatedAt(time.Now().UTC()).
// 		SetUpdatedBy(rankData.UpdatedBy).
// 		Save(ctx)

// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return ierr.WithError(err).
// 				WithHintf("Rank with ID %s was not found", rankData.ID).
// 				WithReportableDetails(map[string]any{
// 					"rank_id": rankData.ID,
// 				}).
// 				Mark(ierr.ErrNotFound)
// 		}
// 		return ierr.WithError(err).
// 			WithHint("Failed to delete rank").
// 			WithReportableDetails(map[string]any{
// 				"rank_id": rankData.ID,
// 			}).
// 			Mark(ierr.ErrDatabase)
// 	}

// 	return nil
// }

// // RankQuery type alias for better readability
// type RankQuery = *ent.RankQuery

// // RankQueryOptions implements query options for rank queries
// type RankQueryOptions struct {
// 	QueryOptionsHelper
// }

// // Ensure RankQueryOptions implements EntityQueryOptions interface
// var _ EntityQueryOptions[RankQuery, *types.RankFilter] = (*RankQueryOptions)(nil)

// func (o RankQueryOptions) ApplyStatusFilter(query RankQuery, status string) RankQuery {
// 	if status == "" {
// 		return query.Where(rank.StatusNotIn(string(types.StatusDeleted)))
// 	}
// 	return query.Where(rank.Status(status))
// }

// func (o RankQueryOptions) ApplySortFilter(query RankQuery, field string, order string) RankQuery {
// 	field, order = o.ValidateSort(field, order)
// 	fieldName := o.GetFieldName(field)
// 	if order == types.OrderDesc {
// 		return query.Order(ent.Desc(fieldName))
// 	}
// 	return query.Order(ent.Asc(fieldName))
// }

// func (o RankQueryOptions) ApplyPaginationFilter(query RankQuery, limit int, offset int) RankQuery {
// 	limit, offset = o.ValidatePagination(limit, offset)
// 	return query.Offset(offset).Limit(limit)
// }

// func (o RankQueryOptions) GetFieldName(field string) string {
// 	switch field {
// 	case "created_at":
// 		return rank.FieldCreatedAt
// 	case "updated_at":
// 		return rank.FieldUpdatedAt
// 	case "name":
// 		return rank.FieldName
// 	case "code":
// 		return rank.FieldCode
// 	case "hierarchy_level":
// 		return rank.FieldHierarchyLevel
// 	default:
// 		return field
// 	}
// }

// func (o RankQueryOptions) ApplyBaseFilters(
// 	_ context.Context,
// 	query RankQuery,
// 	filter *types.RankFilter,
// ) RankQuery {
// 	if filter == nil {
// 		return query.Where(rank.StatusNotIn(string(types.StatusDeleted)))
// 	}

// 	// Apply status filter
// 	query = o.ApplyStatusFilter(query, filter.GetStatus())

// 	// Apply pagination
// 	if !filter.IsUnlimited() {
// 		query = o.ApplyPaginationFilter(query, filter.GetLimit(), filter.GetOffset())
// 	}

// 	// Apply sorting
// 	query = o.ApplySortFilter(query, filter.GetSort(), filter.GetOrder())

// 	return query
// }

// func (o RankQueryOptions) ApplyEntityQueryOptions(
// 	_ context.Context,
// 	f *types.RankFilter,
// 	query RankQuery,
// ) RankQuery {
// 	if f == nil {
// 		return query
// 	}

// 	// Apply rank IDs filter if specified
// 	if len(f.RankIds) > 0 {
// 		query = query.Where(rank.IDIn(f.RankIds...))
// 	}

// 	// Apply hierarchy level filter if specified
// 	if f.Level != nil {
// 		query = query.Where(rank.HierarchyLevelEQ(*f.Level))
// 	}

// 	// Apply time range filters if specified
// 	if f.TimeRangeFilter != nil {
// 		if f.StartTime != nil {
// 			query = query.Where(rank.CreatedAtGTE(*f.StartTime))
// 		}
// 		if f.EndTime != nil {
// 			query = query.Where(rank.CreatedAtLTE(*f.EndTime))
// 		}
// 	}

// 	// Apply expansion if requested
// 	expand := f.GetExpand()
// 	if expand.Has("users") {
// 		query = query.WithUsers()
// 	}

// 	return query
// }
