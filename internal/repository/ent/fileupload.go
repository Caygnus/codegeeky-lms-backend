package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/fileupload"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	domainFileUpload "github.com/omkar273/codegeeky/internal/fileupload"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
)

type fileUploadRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts FileUploadQueryOptions
}

func NewFileUploadRepository(client postgres.IClient, logger *logger.Logger) domainFileUpload.Repository {
	return &fileUploadRepository{
		client:    client,
		log:       *logger,
		queryOpts: FileUploadQueryOptions{},
	}
}

func (r *fileUploadRepository) Create(ctx context.Context, fileData *domainFileUpload.FileUpload) (*domainFileUpload.FileUpload, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating file upload",
		"file_id", fileData.ID,
		"file_name", fileData.FileName,
		"file_type", fileData.FileType,
	)

	entFileUpload, err := client.FileUpload.Create().
		SetID(fileData.ID).
		SetFileName(fileData.FileName).
		SetFileType(fileData.FileType).
		SetExtension(string(fileData.Extension)).
		SetMimeType(fileData.MimeType).
		SetPublicURL(fileData.PublicURL).
		SetProvider(string(fileData.Provider)).
		SetExternalID(fileData.ExternalID).
		SetSizeBytes(fileData.SizeBytes).
		SetStatus(string(fileData.Status)).
		SetCreatedAt(fileData.CreatedAt).
		SetUpdatedAt(fileData.UpdatedAt).
		SetCreatedBy(fileData.CreatedBy).
		SetUpdatedBy(fileData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ierr.WithError(err).
				WithHint("File upload with this external ID already exists").
				WithReportableDetails(map[string]any{
					"file_id":     fileData.ID,
					"external_id": fileData.ExternalID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to create file upload").
			WithReportableDetails(map[string]any{
				"file_id":   fileData.ID,
				"file_name": fileData.FileName,
			}).
			Mark(ierr.ErrDatabase)
	}

	// Set optional fields if they exist
	if fileData.SecureURL != nil {
		entFileUpload, err = client.FileUpload.UpdateOneID(entFileUpload.ID).
			SetSecureURL(*fileData.SecureURL).
			Save(ctx)
		if err != nil {
			return nil, ierr.WithError(err).
				WithHint("Failed to update file upload with secure URL").
				WithReportableDetails(map[string]any{
					"file_id": fileData.ID,
				}).
				Mark(ierr.ErrDatabase)
		}
	}

	if fileData.FileSize != nil {
		entFileUpload, err = client.FileUpload.UpdateOneID(entFileUpload.ID).
			SetFileSize(*fileData.FileSize).
			Save(ctx)
		if err != nil {
			return nil, ierr.WithError(err).
				WithHint("Failed to update file upload with file size").
				WithReportableDetails(map[string]any{
					"file_id": fileData.ID,
				}).
				Mark(ierr.ErrDatabase)
		}
	}

	return domainFileUpload.FromEnt(entFileUpload), nil
}

func (r *fileUploadRepository) Get(ctx context.Context, id string) (*domainFileUpload.FileUpload, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting file upload", "file_id", id)

	entFileUpload, err := client.FileUpload.Query().
		Where(fileupload.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("File upload with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"file_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get file upload").
			WithReportableDetails(map[string]any{
				"file_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainFileUpload.FromEnt(entFileUpload), nil
}

func (r *fileUploadRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting file upload", "file_id", id)

	_, err := client.FileUpload.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("File upload with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"file_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete file upload").
			WithReportableDetails(map[string]any{
				"file_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *fileUploadRepository) Exists(ctx context.Context, id string) (bool, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("checking if file upload exists", "file_id", id)

	exists, err := client.FileUpload.Query().
		Where(
			fileupload.ID(id),
			fileupload.StatusNotIn(string(types.StatusDeleted)),
		).
		Exist(ctx)

	if err != nil {
		return false, ierr.WithError(err).
			WithHint("Failed to check if file upload exists").
			WithReportableDetails(map[string]any{
				"file_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return exists, nil
}

// FileUploadQuery type alias for better readability
type FileUploadQuery = *ent.FileUploadQuery

// FileUploadQueryOptions implements query options for file upload queries
type FileUploadQueryOptions struct {
	QueryOptionsHelper
}

// Ensure FileUploadQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[FileUploadQuery, *types.FileUploadFilter] = (*FileUploadQueryOptions)(nil)

func (o FileUploadQueryOptions) ApplyStatusFilter(query FileUploadQuery, status string) FileUploadQuery {
	if status == "" {
		return query.Where(fileupload.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(fileupload.Status(status))
}

func (o FileUploadQueryOptions) ApplySortFilter(query FileUploadQuery, field string, order string) FileUploadQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o FileUploadQueryOptions) ApplyPaginationFilter(query FileUploadQuery, limit int, offset int) FileUploadQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o FileUploadQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return fileupload.FieldCreatedAt
	case "updated_at":
		return fileupload.FieldUpdatedAt
	case "file_name":
		return fileupload.FieldFileName
	case "file_type":
		return fileupload.FieldFileType
	case "extension":
		return fileupload.FieldExtension
	case "mime_type":
		return fileupload.FieldMimeType
	case "public_url":
		return fileupload.FieldPublicURL
	case "secure_url":
		return fileupload.FieldSecureURL
	case "provider":
		return fileupload.FieldProvider
	case "external_id":
		return fileupload.FieldExternalID
	case "size_bytes":
		return fileupload.FieldSizeBytes
	case "file_size":
		return fileupload.FieldFileSize
	case "created_by":
		return fileupload.FieldCreatedBy
	case "updated_by":
		return fileupload.FieldUpdatedBy
	case "status":
		return fileupload.FieldStatus
	case "id":
		return fileupload.FieldID
	default:
		return field
	}
}

func (o FileUploadQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query FileUploadQuery,
	filter *types.FileUploadFilter,
) FileUploadQuery {
	if filter == nil {
		return query.Where(fileupload.StatusNotIn(string(types.StatusDeleted)))
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

func (o FileUploadQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.FileUploadFilter,
	query FileUploadQuery,
) FileUploadQuery {
	if f == nil {
		return query
	}

	// Apply external IDs filter if specified
	if len(f.ExternalIDs) > 0 {
		query = query.Where(fileupload.ExternalIDIn(f.ExternalIDs...))
	}

	// Apply file IDs filter if specified
	if len(f.FileIDs) > 0 {
		query = query.Where(fileupload.IDIn(f.FileIDs...))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(fileupload.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(fileupload.CreatedAtLTE(*f.EndTime))
		}
	}

	return query
}
