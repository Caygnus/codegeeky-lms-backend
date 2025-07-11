// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/omkar273/codegeeky/ent/fileupload"
)

// FileUploadCreate is the builder for creating a FileUpload entity.
type FileUploadCreate struct {
	config
	mutation *FileUploadMutation
	hooks    []Hook
}

// SetStatus sets the "status" field.
func (fuc *FileUploadCreate) SetStatus(s string) *FileUploadCreate {
	fuc.mutation.SetStatus(s)
	return fuc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableStatus(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetStatus(*s)
	}
	return fuc
}

// SetCreatedAt sets the "created_at" field.
func (fuc *FileUploadCreate) SetCreatedAt(t time.Time) *FileUploadCreate {
	fuc.mutation.SetCreatedAt(t)
	return fuc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableCreatedAt(t *time.Time) *FileUploadCreate {
	if t != nil {
		fuc.SetCreatedAt(*t)
	}
	return fuc
}

// SetUpdatedAt sets the "updated_at" field.
func (fuc *FileUploadCreate) SetUpdatedAt(t time.Time) *FileUploadCreate {
	fuc.mutation.SetUpdatedAt(t)
	return fuc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableUpdatedAt(t *time.Time) *FileUploadCreate {
	if t != nil {
		fuc.SetUpdatedAt(*t)
	}
	return fuc
}

// SetCreatedBy sets the "created_by" field.
func (fuc *FileUploadCreate) SetCreatedBy(s string) *FileUploadCreate {
	fuc.mutation.SetCreatedBy(s)
	return fuc
}

// SetNillableCreatedBy sets the "created_by" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableCreatedBy(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetCreatedBy(*s)
	}
	return fuc
}

// SetUpdatedBy sets the "updated_by" field.
func (fuc *FileUploadCreate) SetUpdatedBy(s string) *FileUploadCreate {
	fuc.mutation.SetUpdatedBy(s)
	return fuc
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableUpdatedBy(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetUpdatedBy(*s)
	}
	return fuc
}

// SetFileName sets the "file_name" field.
func (fuc *FileUploadCreate) SetFileName(s string) *FileUploadCreate {
	fuc.mutation.SetFileName(s)
	return fuc
}

// SetFileType sets the "file_type" field.
func (fuc *FileUploadCreate) SetFileType(s string) *FileUploadCreate {
	fuc.mutation.SetFileType(s)
	return fuc
}

// SetExtension sets the "extension" field.
func (fuc *FileUploadCreate) SetExtension(s string) *FileUploadCreate {
	fuc.mutation.SetExtension(s)
	return fuc
}

// SetMimeType sets the "mime_type" field.
func (fuc *FileUploadCreate) SetMimeType(s string) *FileUploadCreate {
	fuc.mutation.SetMimeType(s)
	return fuc
}

// SetPublicURL sets the "public_url" field.
func (fuc *FileUploadCreate) SetPublicURL(s string) *FileUploadCreate {
	fuc.mutation.SetPublicURL(s)
	return fuc
}

// SetSecureURL sets the "secure_url" field.
func (fuc *FileUploadCreate) SetSecureURL(s string) *FileUploadCreate {
	fuc.mutation.SetSecureURL(s)
	return fuc
}

// SetNillableSecureURL sets the "secure_url" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableSecureURL(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetSecureURL(*s)
	}
	return fuc
}

// SetProvider sets the "provider" field.
func (fuc *FileUploadCreate) SetProvider(s string) *FileUploadCreate {
	fuc.mutation.SetProvider(s)
	return fuc
}

// SetExternalID sets the "external_id" field.
func (fuc *FileUploadCreate) SetExternalID(s string) *FileUploadCreate {
	fuc.mutation.SetExternalID(s)
	return fuc
}

// SetSizeBytes sets the "size_bytes" field.
func (fuc *FileUploadCreate) SetSizeBytes(i int64) *FileUploadCreate {
	fuc.mutation.SetSizeBytes(i)
	return fuc
}

// SetFileSize sets the "file_size" field.
func (fuc *FileUploadCreate) SetFileSize(s string) *FileUploadCreate {
	fuc.mutation.SetFileSize(s)
	return fuc
}

// SetNillableFileSize sets the "file_size" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableFileSize(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetFileSize(*s)
	}
	return fuc
}

// SetID sets the "id" field.
func (fuc *FileUploadCreate) SetID(s string) *FileUploadCreate {
	fuc.mutation.SetID(s)
	return fuc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (fuc *FileUploadCreate) SetNillableID(s *string) *FileUploadCreate {
	if s != nil {
		fuc.SetID(*s)
	}
	return fuc
}

// Mutation returns the FileUploadMutation object of the builder.
func (fuc *FileUploadCreate) Mutation() *FileUploadMutation {
	return fuc.mutation
}

// Save creates the FileUpload in the database.
func (fuc *FileUploadCreate) Save(ctx context.Context) (*FileUpload, error) {
	fuc.defaults()
	return withHooks(ctx, fuc.sqlSave, fuc.mutation, fuc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fuc *FileUploadCreate) SaveX(ctx context.Context) *FileUpload {
	v, err := fuc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fuc *FileUploadCreate) Exec(ctx context.Context) error {
	_, err := fuc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuc *FileUploadCreate) ExecX(ctx context.Context) {
	if err := fuc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fuc *FileUploadCreate) defaults() {
	if _, ok := fuc.mutation.Status(); !ok {
		v := fileupload.DefaultStatus
		fuc.mutation.SetStatus(v)
	}
	if _, ok := fuc.mutation.CreatedAt(); !ok {
		v := fileupload.DefaultCreatedAt()
		fuc.mutation.SetCreatedAt(v)
	}
	if _, ok := fuc.mutation.UpdatedAt(); !ok {
		v := fileupload.DefaultUpdatedAt()
		fuc.mutation.SetUpdatedAt(v)
	}
	if _, ok := fuc.mutation.ID(); !ok {
		v := fileupload.DefaultID()
		fuc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fuc *FileUploadCreate) check() error {
	if _, ok := fuc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "FileUpload.status"`)}
	}
	if _, ok := fuc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "FileUpload.created_at"`)}
	}
	if _, ok := fuc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "FileUpload.updated_at"`)}
	}
	if _, ok := fuc.mutation.FileName(); !ok {
		return &ValidationError{Name: "file_name", err: errors.New(`ent: missing required field "FileUpload.file_name"`)}
	}
	if v, ok := fuc.mutation.FileName(); ok {
		if err := fileupload.FileNameValidator(v); err != nil {
			return &ValidationError{Name: "file_name", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_name": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.FileType(); !ok {
		return &ValidationError{Name: "file_type", err: errors.New(`ent: missing required field "FileUpload.file_type"`)}
	}
	if v, ok := fuc.mutation.FileType(); ok {
		if err := fileupload.FileTypeValidator(v); err != nil {
			return &ValidationError{Name: "file_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_type": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.Extension(); !ok {
		return &ValidationError{Name: "extension", err: errors.New(`ent: missing required field "FileUpload.extension"`)}
	}
	if v, ok := fuc.mutation.Extension(); ok {
		if err := fileupload.ExtensionValidator(v); err != nil {
			return &ValidationError{Name: "extension", err: fmt.Errorf(`ent: validator failed for field "FileUpload.extension": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.MimeType(); !ok {
		return &ValidationError{Name: "mime_type", err: errors.New(`ent: missing required field "FileUpload.mime_type"`)}
	}
	if v, ok := fuc.mutation.MimeType(); ok {
		if err := fileupload.MimeTypeValidator(v); err != nil {
			return &ValidationError{Name: "mime_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.mime_type": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.PublicURL(); !ok {
		return &ValidationError{Name: "public_url", err: errors.New(`ent: missing required field "FileUpload.public_url"`)}
	}
	if v, ok := fuc.mutation.PublicURL(); ok {
		if err := fileupload.PublicURLValidator(v); err != nil {
			return &ValidationError{Name: "public_url", err: fmt.Errorf(`ent: validator failed for field "FileUpload.public_url": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.Provider(); !ok {
		return &ValidationError{Name: "provider", err: errors.New(`ent: missing required field "FileUpload.provider"`)}
	}
	if v, ok := fuc.mutation.Provider(); ok {
		if err := fileupload.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "FileUpload.provider": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.ExternalID(); !ok {
		return &ValidationError{Name: "external_id", err: errors.New(`ent: missing required field "FileUpload.external_id"`)}
	}
	if v, ok := fuc.mutation.ExternalID(); ok {
		if err := fileupload.ExternalIDValidator(v); err != nil {
			return &ValidationError{Name: "external_id", err: fmt.Errorf(`ent: validator failed for field "FileUpload.external_id": %w`, err)}
		}
	}
	if _, ok := fuc.mutation.SizeBytes(); !ok {
		return &ValidationError{Name: "size_bytes", err: errors.New(`ent: missing required field "FileUpload.size_bytes"`)}
	}
	if v, ok := fuc.mutation.SizeBytes(); ok {
		if err := fileupload.SizeBytesValidator(v); err != nil {
			return &ValidationError{Name: "size_bytes", err: fmt.Errorf(`ent: validator failed for field "FileUpload.size_bytes": %w`, err)}
		}
	}
	return nil
}

func (fuc *FileUploadCreate) sqlSave(ctx context.Context) (*FileUpload, error) {
	if err := fuc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fuc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fuc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected FileUpload.ID type: %T", _spec.ID.Value)
		}
	}
	fuc.mutation.id = &_node.ID
	fuc.mutation.done = true
	return _node, nil
}

func (fuc *FileUploadCreate) createSpec() (*FileUpload, *sqlgraph.CreateSpec) {
	var (
		_node = &FileUpload{config: fuc.config}
		_spec = sqlgraph.NewCreateSpec(fileupload.Table, sqlgraph.NewFieldSpec(fileupload.FieldID, field.TypeString))
	)
	if id, ok := fuc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := fuc.mutation.Status(); ok {
		_spec.SetField(fileupload.FieldStatus, field.TypeString, value)
		_node.Status = value
	}
	if value, ok := fuc.mutation.CreatedAt(); ok {
		_spec.SetField(fileupload.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := fuc.mutation.UpdatedAt(); ok {
		_spec.SetField(fileupload.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := fuc.mutation.CreatedBy(); ok {
		_spec.SetField(fileupload.FieldCreatedBy, field.TypeString, value)
		_node.CreatedBy = value
	}
	if value, ok := fuc.mutation.UpdatedBy(); ok {
		_spec.SetField(fileupload.FieldUpdatedBy, field.TypeString, value)
		_node.UpdatedBy = value
	}
	if value, ok := fuc.mutation.FileName(); ok {
		_spec.SetField(fileupload.FieldFileName, field.TypeString, value)
		_node.FileName = value
	}
	if value, ok := fuc.mutation.FileType(); ok {
		_spec.SetField(fileupload.FieldFileType, field.TypeString, value)
		_node.FileType = value
	}
	if value, ok := fuc.mutation.Extension(); ok {
		_spec.SetField(fileupload.FieldExtension, field.TypeString, value)
		_node.Extension = value
	}
	if value, ok := fuc.mutation.MimeType(); ok {
		_spec.SetField(fileupload.FieldMimeType, field.TypeString, value)
		_node.MimeType = value
	}
	if value, ok := fuc.mutation.PublicURL(); ok {
		_spec.SetField(fileupload.FieldPublicURL, field.TypeString, value)
		_node.PublicURL = value
	}
	if value, ok := fuc.mutation.SecureURL(); ok {
		_spec.SetField(fileupload.FieldSecureURL, field.TypeString, value)
		_node.SecureURL = &value
	}
	if value, ok := fuc.mutation.Provider(); ok {
		_spec.SetField(fileupload.FieldProvider, field.TypeString, value)
		_node.Provider = value
	}
	if value, ok := fuc.mutation.ExternalID(); ok {
		_spec.SetField(fileupload.FieldExternalID, field.TypeString, value)
		_node.ExternalID = value
	}
	if value, ok := fuc.mutation.SizeBytes(); ok {
		_spec.SetField(fileupload.FieldSizeBytes, field.TypeInt64, value)
		_node.SizeBytes = value
	}
	if value, ok := fuc.mutation.FileSize(); ok {
		_spec.SetField(fileupload.FieldFileSize, field.TypeString, value)
		_node.FileSize = &value
	}
	return _node, _spec
}

// FileUploadCreateBulk is the builder for creating many FileUpload entities in bulk.
type FileUploadCreateBulk struct {
	config
	err      error
	builders []*FileUploadCreate
}

// Save creates the FileUpload entities in the database.
func (fucb *FileUploadCreateBulk) Save(ctx context.Context) ([]*FileUpload, error) {
	if fucb.err != nil {
		return nil, fucb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(fucb.builders))
	nodes := make([]*FileUpload, len(fucb.builders))
	mutators := make([]Mutator, len(fucb.builders))
	for i := range fucb.builders {
		func(i int, root context.Context) {
			builder := fucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FileUploadMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fucb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fucb *FileUploadCreateBulk) SaveX(ctx context.Context) []*FileUpload {
	v, err := fucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fucb *FileUploadCreateBulk) Exec(ctx context.Context) error {
	_, err := fucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fucb *FileUploadCreateBulk) ExecX(ctx context.Context) {
	if err := fucb.Exec(ctx); err != nil {
		panic(err)
	}
}
