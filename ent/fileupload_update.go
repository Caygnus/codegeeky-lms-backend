// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/omkar273/codegeeky/ent/fileupload"
	"github.com/omkar273/codegeeky/ent/predicate"
)

// FileUploadUpdate is the builder for updating FileUpload entities.
type FileUploadUpdate struct {
	config
	hooks    []Hook
	mutation *FileUploadMutation
}

// Where appends a list predicates to the FileUploadUpdate builder.
func (fuu *FileUploadUpdate) Where(ps ...predicate.FileUpload) *FileUploadUpdate {
	fuu.mutation.Where(ps...)
	return fuu
}

// SetStatus sets the "status" field.
func (fuu *FileUploadUpdate) SetStatus(s string) *FileUploadUpdate {
	fuu.mutation.SetStatus(s)
	return fuu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableStatus(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetStatus(*s)
	}
	return fuu
}

// SetUpdatedAt sets the "updated_at" field.
func (fuu *FileUploadUpdate) SetUpdatedAt(t time.Time) *FileUploadUpdate {
	fuu.mutation.SetUpdatedAt(t)
	return fuu
}

// SetUpdatedBy sets the "updated_by" field.
func (fuu *FileUploadUpdate) SetUpdatedBy(s string) *FileUploadUpdate {
	fuu.mutation.SetUpdatedBy(s)
	return fuu
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableUpdatedBy(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetUpdatedBy(*s)
	}
	return fuu
}

// ClearUpdatedBy clears the value of the "updated_by" field.
func (fuu *FileUploadUpdate) ClearUpdatedBy() *FileUploadUpdate {
	fuu.mutation.ClearUpdatedBy()
	return fuu
}

// SetFileName sets the "file_name" field.
func (fuu *FileUploadUpdate) SetFileName(s string) *FileUploadUpdate {
	fuu.mutation.SetFileName(s)
	return fuu
}

// SetNillableFileName sets the "file_name" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableFileName(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetFileName(*s)
	}
	return fuu
}

// SetFileType sets the "file_type" field.
func (fuu *FileUploadUpdate) SetFileType(s string) *FileUploadUpdate {
	fuu.mutation.SetFileType(s)
	return fuu
}

// SetNillableFileType sets the "file_type" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableFileType(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetFileType(*s)
	}
	return fuu
}

// SetExtension sets the "extension" field.
func (fuu *FileUploadUpdate) SetExtension(s string) *FileUploadUpdate {
	fuu.mutation.SetExtension(s)
	return fuu
}

// SetNillableExtension sets the "extension" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableExtension(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetExtension(*s)
	}
	return fuu
}

// SetMimeType sets the "mime_type" field.
func (fuu *FileUploadUpdate) SetMimeType(s string) *FileUploadUpdate {
	fuu.mutation.SetMimeType(s)
	return fuu
}

// SetNillableMimeType sets the "mime_type" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableMimeType(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetMimeType(*s)
	}
	return fuu
}

// SetPublicURL sets the "public_url" field.
func (fuu *FileUploadUpdate) SetPublicURL(s string) *FileUploadUpdate {
	fuu.mutation.SetPublicURL(s)
	return fuu
}

// SetNillablePublicURL sets the "public_url" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillablePublicURL(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetPublicURL(*s)
	}
	return fuu
}

// SetSecureURL sets the "secure_url" field.
func (fuu *FileUploadUpdate) SetSecureURL(s string) *FileUploadUpdate {
	fuu.mutation.SetSecureURL(s)
	return fuu
}

// SetNillableSecureURL sets the "secure_url" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableSecureURL(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetSecureURL(*s)
	}
	return fuu
}

// ClearSecureURL clears the value of the "secure_url" field.
func (fuu *FileUploadUpdate) ClearSecureURL() *FileUploadUpdate {
	fuu.mutation.ClearSecureURL()
	return fuu
}

// SetProvider sets the "provider" field.
func (fuu *FileUploadUpdate) SetProvider(s string) *FileUploadUpdate {
	fuu.mutation.SetProvider(s)
	return fuu
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableProvider(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetProvider(*s)
	}
	return fuu
}

// SetExternalID sets the "external_id" field.
func (fuu *FileUploadUpdate) SetExternalID(s string) *FileUploadUpdate {
	fuu.mutation.SetExternalID(s)
	return fuu
}

// SetNillableExternalID sets the "external_id" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableExternalID(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetExternalID(*s)
	}
	return fuu
}

// SetSizeBytes sets the "size_bytes" field.
func (fuu *FileUploadUpdate) SetSizeBytes(i int64) *FileUploadUpdate {
	fuu.mutation.ResetSizeBytes()
	fuu.mutation.SetSizeBytes(i)
	return fuu
}

// SetNillableSizeBytes sets the "size_bytes" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableSizeBytes(i *int64) *FileUploadUpdate {
	if i != nil {
		fuu.SetSizeBytes(*i)
	}
	return fuu
}

// AddSizeBytes adds i to the "size_bytes" field.
func (fuu *FileUploadUpdate) AddSizeBytes(i int64) *FileUploadUpdate {
	fuu.mutation.AddSizeBytes(i)
	return fuu
}

// SetFileSize sets the "file_size" field.
func (fuu *FileUploadUpdate) SetFileSize(s string) *FileUploadUpdate {
	fuu.mutation.SetFileSize(s)
	return fuu
}

// SetNillableFileSize sets the "file_size" field if the given value is not nil.
func (fuu *FileUploadUpdate) SetNillableFileSize(s *string) *FileUploadUpdate {
	if s != nil {
		fuu.SetFileSize(*s)
	}
	return fuu
}

// ClearFileSize clears the value of the "file_size" field.
func (fuu *FileUploadUpdate) ClearFileSize() *FileUploadUpdate {
	fuu.mutation.ClearFileSize()
	return fuu
}

// Mutation returns the FileUploadMutation object of the builder.
func (fuu *FileUploadUpdate) Mutation() *FileUploadMutation {
	return fuu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fuu *FileUploadUpdate) Save(ctx context.Context) (int, error) {
	fuu.defaults()
	return withHooks(ctx, fuu.sqlSave, fuu.mutation, fuu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fuu *FileUploadUpdate) SaveX(ctx context.Context) int {
	affected, err := fuu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fuu *FileUploadUpdate) Exec(ctx context.Context) error {
	_, err := fuu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuu *FileUploadUpdate) ExecX(ctx context.Context) {
	if err := fuu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fuu *FileUploadUpdate) defaults() {
	if _, ok := fuu.mutation.UpdatedAt(); !ok {
		v := fileupload.UpdateDefaultUpdatedAt()
		fuu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fuu *FileUploadUpdate) check() error {
	if v, ok := fuu.mutation.FileName(); ok {
		if err := fileupload.FileNameValidator(v); err != nil {
			return &ValidationError{Name: "file_name", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_name": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.FileType(); ok {
		if err := fileupload.FileTypeValidator(v); err != nil {
			return &ValidationError{Name: "file_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_type": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.Extension(); ok {
		if err := fileupload.ExtensionValidator(v); err != nil {
			return &ValidationError{Name: "extension", err: fmt.Errorf(`ent: validator failed for field "FileUpload.extension": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.MimeType(); ok {
		if err := fileupload.MimeTypeValidator(v); err != nil {
			return &ValidationError{Name: "mime_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.mime_type": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.PublicURL(); ok {
		if err := fileupload.PublicURLValidator(v); err != nil {
			return &ValidationError{Name: "public_url", err: fmt.Errorf(`ent: validator failed for field "FileUpload.public_url": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.Provider(); ok {
		if err := fileupload.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "FileUpload.provider": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.ExternalID(); ok {
		if err := fileupload.ExternalIDValidator(v); err != nil {
			return &ValidationError{Name: "external_id", err: fmt.Errorf(`ent: validator failed for field "FileUpload.external_id": %w`, err)}
		}
	}
	if v, ok := fuu.mutation.SizeBytes(); ok {
		if err := fileupload.SizeBytesValidator(v); err != nil {
			return &ValidationError{Name: "size_bytes", err: fmt.Errorf(`ent: validator failed for field "FileUpload.size_bytes": %w`, err)}
		}
	}
	return nil
}

func (fuu *FileUploadUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := fuu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(fileupload.Table, fileupload.Columns, sqlgraph.NewFieldSpec(fileupload.FieldID, field.TypeString))
	if ps := fuu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuu.mutation.Status(); ok {
		_spec.SetField(fileupload.FieldStatus, field.TypeString, value)
	}
	if value, ok := fuu.mutation.UpdatedAt(); ok {
		_spec.SetField(fileupload.FieldUpdatedAt, field.TypeTime, value)
	}
	if fuu.mutation.CreatedByCleared() {
		_spec.ClearField(fileupload.FieldCreatedBy, field.TypeString)
	}
	if value, ok := fuu.mutation.UpdatedBy(); ok {
		_spec.SetField(fileupload.FieldUpdatedBy, field.TypeString, value)
	}
	if fuu.mutation.UpdatedByCleared() {
		_spec.ClearField(fileupload.FieldUpdatedBy, field.TypeString)
	}
	if value, ok := fuu.mutation.FileName(); ok {
		_spec.SetField(fileupload.FieldFileName, field.TypeString, value)
	}
	if value, ok := fuu.mutation.FileType(); ok {
		_spec.SetField(fileupload.FieldFileType, field.TypeString, value)
	}
	if value, ok := fuu.mutation.Extension(); ok {
		_spec.SetField(fileupload.FieldExtension, field.TypeString, value)
	}
	if value, ok := fuu.mutation.MimeType(); ok {
		_spec.SetField(fileupload.FieldMimeType, field.TypeString, value)
	}
	if value, ok := fuu.mutation.PublicURL(); ok {
		_spec.SetField(fileupload.FieldPublicURL, field.TypeString, value)
	}
	if value, ok := fuu.mutation.SecureURL(); ok {
		_spec.SetField(fileupload.FieldSecureURL, field.TypeString, value)
	}
	if fuu.mutation.SecureURLCleared() {
		_spec.ClearField(fileupload.FieldSecureURL, field.TypeString)
	}
	if value, ok := fuu.mutation.Provider(); ok {
		_spec.SetField(fileupload.FieldProvider, field.TypeString, value)
	}
	if value, ok := fuu.mutation.ExternalID(); ok {
		_spec.SetField(fileupload.FieldExternalID, field.TypeString, value)
	}
	if value, ok := fuu.mutation.SizeBytes(); ok {
		_spec.SetField(fileupload.FieldSizeBytes, field.TypeInt64, value)
	}
	if value, ok := fuu.mutation.AddedSizeBytes(); ok {
		_spec.AddField(fileupload.FieldSizeBytes, field.TypeInt64, value)
	}
	if value, ok := fuu.mutation.FileSize(); ok {
		_spec.SetField(fileupload.FieldFileSize, field.TypeString, value)
	}
	if fuu.mutation.FileSizeCleared() {
		_spec.ClearField(fileupload.FieldFileSize, field.TypeString)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fuu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fileupload.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fuu.mutation.done = true
	return n, nil
}

// FileUploadUpdateOne is the builder for updating a single FileUpload entity.
type FileUploadUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FileUploadMutation
}

// SetStatus sets the "status" field.
func (fuuo *FileUploadUpdateOne) SetStatus(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetStatus(s)
	return fuuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableStatus(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetStatus(*s)
	}
	return fuuo
}

// SetUpdatedAt sets the "updated_at" field.
func (fuuo *FileUploadUpdateOne) SetUpdatedAt(t time.Time) *FileUploadUpdateOne {
	fuuo.mutation.SetUpdatedAt(t)
	return fuuo
}

// SetUpdatedBy sets the "updated_by" field.
func (fuuo *FileUploadUpdateOne) SetUpdatedBy(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetUpdatedBy(s)
	return fuuo
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableUpdatedBy(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetUpdatedBy(*s)
	}
	return fuuo
}

// ClearUpdatedBy clears the value of the "updated_by" field.
func (fuuo *FileUploadUpdateOne) ClearUpdatedBy() *FileUploadUpdateOne {
	fuuo.mutation.ClearUpdatedBy()
	return fuuo
}

// SetFileName sets the "file_name" field.
func (fuuo *FileUploadUpdateOne) SetFileName(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetFileName(s)
	return fuuo
}

// SetNillableFileName sets the "file_name" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableFileName(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetFileName(*s)
	}
	return fuuo
}

// SetFileType sets the "file_type" field.
func (fuuo *FileUploadUpdateOne) SetFileType(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetFileType(s)
	return fuuo
}

// SetNillableFileType sets the "file_type" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableFileType(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetFileType(*s)
	}
	return fuuo
}

// SetExtension sets the "extension" field.
func (fuuo *FileUploadUpdateOne) SetExtension(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetExtension(s)
	return fuuo
}

// SetNillableExtension sets the "extension" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableExtension(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetExtension(*s)
	}
	return fuuo
}

// SetMimeType sets the "mime_type" field.
func (fuuo *FileUploadUpdateOne) SetMimeType(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetMimeType(s)
	return fuuo
}

// SetNillableMimeType sets the "mime_type" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableMimeType(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetMimeType(*s)
	}
	return fuuo
}

// SetPublicURL sets the "public_url" field.
func (fuuo *FileUploadUpdateOne) SetPublicURL(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetPublicURL(s)
	return fuuo
}

// SetNillablePublicURL sets the "public_url" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillablePublicURL(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetPublicURL(*s)
	}
	return fuuo
}

// SetSecureURL sets the "secure_url" field.
func (fuuo *FileUploadUpdateOne) SetSecureURL(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetSecureURL(s)
	return fuuo
}

// SetNillableSecureURL sets the "secure_url" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableSecureURL(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetSecureURL(*s)
	}
	return fuuo
}

// ClearSecureURL clears the value of the "secure_url" field.
func (fuuo *FileUploadUpdateOne) ClearSecureURL() *FileUploadUpdateOne {
	fuuo.mutation.ClearSecureURL()
	return fuuo
}

// SetProvider sets the "provider" field.
func (fuuo *FileUploadUpdateOne) SetProvider(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetProvider(s)
	return fuuo
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableProvider(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetProvider(*s)
	}
	return fuuo
}

// SetExternalID sets the "external_id" field.
func (fuuo *FileUploadUpdateOne) SetExternalID(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetExternalID(s)
	return fuuo
}

// SetNillableExternalID sets the "external_id" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableExternalID(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetExternalID(*s)
	}
	return fuuo
}

// SetSizeBytes sets the "size_bytes" field.
func (fuuo *FileUploadUpdateOne) SetSizeBytes(i int64) *FileUploadUpdateOne {
	fuuo.mutation.ResetSizeBytes()
	fuuo.mutation.SetSizeBytes(i)
	return fuuo
}

// SetNillableSizeBytes sets the "size_bytes" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableSizeBytes(i *int64) *FileUploadUpdateOne {
	if i != nil {
		fuuo.SetSizeBytes(*i)
	}
	return fuuo
}

// AddSizeBytes adds i to the "size_bytes" field.
func (fuuo *FileUploadUpdateOne) AddSizeBytes(i int64) *FileUploadUpdateOne {
	fuuo.mutation.AddSizeBytes(i)
	return fuuo
}

// SetFileSize sets the "file_size" field.
func (fuuo *FileUploadUpdateOne) SetFileSize(s string) *FileUploadUpdateOne {
	fuuo.mutation.SetFileSize(s)
	return fuuo
}

// SetNillableFileSize sets the "file_size" field if the given value is not nil.
func (fuuo *FileUploadUpdateOne) SetNillableFileSize(s *string) *FileUploadUpdateOne {
	if s != nil {
		fuuo.SetFileSize(*s)
	}
	return fuuo
}

// ClearFileSize clears the value of the "file_size" field.
func (fuuo *FileUploadUpdateOne) ClearFileSize() *FileUploadUpdateOne {
	fuuo.mutation.ClearFileSize()
	return fuuo
}

// Mutation returns the FileUploadMutation object of the builder.
func (fuuo *FileUploadUpdateOne) Mutation() *FileUploadMutation {
	return fuuo.mutation
}

// Where appends a list predicates to the FileUploadUpdate builder.
func (fuuo *FileUploadUpdateOne) Where(ps ...predicate.FileUpload) *FileUploadUpdateOne {
	fuuo.mutation.Where(ps...)
	return fuuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuuo *FileUploadUpdateOne) Select(field string, fields ...string) *FileUploadUpdateOne {
	fuuo.fields = append([]string{field}, fields...)
	return fuuo
}

// Save executes the query and returns the updated FileUpload entity.
func (fuuo *FileUploadUpdateOne) Save(ctx context.Context) (*FileUpload, error) {
	fuuo.defaults()
	return withHooks(ctx, fuuo.sqlSave, fuuo.mutation, fuuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fuuo *FileUploadUpdateOne) SaveX(ctx context.Context) *FileUpload {
	node, err := fuuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuuo *FileUploadUpdateOne) Exec(ctx context.Context) error {
	_, err := fuuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuuo *FileUploadUpdateOne) ExecX(ctx context.Context) {
	if err := fuuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fuuo *FileUploadUpdateOne) defaults() {
	if _, ok := fuuo.mutation.UpdatedAt(); !ok {
		v := fileupload.UpdateDefaultUpdatedAt()
		fuuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fuuo *FileUploadUpdateOne) check() error {
	if v, ok := fuuo.mutation.FileName(); ok {
		if err := fileupload.FileNameValidator(v); err != nil {
			return &ValidationError{Name: "file_name", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_name": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.FileType(); ok {
		if err := fileupload.FileTypeValidator(v); err != nil {
			return &ValidationError{Name: "file_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.file_type": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.Extension(); ok {
		if err := fileupload.ExtensionValidator(v); err != nil {
			return &ValidationError{Name: "extension", err: fmt.Errorf(`ent: validator failed for field "FileUpload.extension": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.MimeType(); ok {
		if err := fileupload.MimeTypeValidator(v); err != nil {
			return &ValidationError{Name: "mime_type", err: fmt.Errorf(`ent: validator failed for field "FileUpload.mime_type": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.PublicURL(); ok {
		if err := fileupload.PublicURLValidator(v); err != nil {
			return &ValidationError{Name: "public_url", err: fmt.Errorf(`ent: validator failed for field "FileUpload.public_url": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.Provider(); ok {
		if err := fileupload.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "FileUpload.provider": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.ExternalID(); ok {
		if err := fileupload.ExternalIDValidator(v); err != nil {
			return &ValidationError{Name: "external_id", err: fmt.Errorf(`ent: validator failed for field "FileUpload.external_id": %w`, err)}
		}
	}
	if v, ok := fuuo.mutation.SizeBytes(); ok {
		if err := fileupload.SizeBytesValidator(v); err != nil {
			return &ValidationError{Name: "size_bytes", err: fmt.Errorf(`ent: validator failed for field "FileUpload.size_bytes": %w`, err)}
		}
	}
	return nil
}

func (fuuo *FileUploadUpdateOne) sqlSave(ctx context.Context) (_node *FileUpload, err error) {
	if err := fuuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(fileupload.Table, fileupload.Columns, sqlgraph.NewFieldSpec(fileupload.FieldID, field.TypeString))
	id, ok := fuuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "FileUpload.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, fileupload.FieldID)
		for _, f := range fields {
			if !fileupload.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != fileupload.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fuuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuuo.mutation.Status(); ok {
		_spec.SetField(fileupload.FieldStatus, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.UpdatedAt(); ok {
		_spec.SetField(fileupload.FieldUpdatedAt, field.TypeTime, value)
	}
	if fuuo.mutation.CreatedByCleared() {
		_spec.ClearField(fileupload.FieldCreatedBy, field.TypeString)
	}
	if value, ok := fuuo.mutation.UpdatedBy(); ok {
		_spec.SetField(fileupload.FieldUpdatedBy, field.TypeString, value)
	}
	if fuuo.mutation.UpdatedByCleared() {
		_spec.ClearField(fileupload.FieldUpdatedBy, field.TypeString)
	}
	if value, ok := fuuo.mutation.FileName(); ok {
		_spec.SetField(fileupload.FieldFileName, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.FileType(); ok {
		_spec.SetField(fileupload.FieldFileType, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.Extension(); ok {
		_spec.SetField(fileupload.FieldExtension, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.MimeType(); ok {
		_spec.SetField(fileupload.FieldMimeType, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.PublicURL(); ok {
		_spec.SetField(fileupload.FieldPublicURL, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.SecureURL(); ok {
		_spec.SetField(fileupload.FieldSecureURL, field.TypeString, value)
	}
	if fuuo.mutation.SecureURLCleared() {
		_spec.ClearField(fileupload.FieldSecureURL, field.TypeString)
	}
	if value, ok := fuuo.mutation.Provider(); ok {
		_spec.SetField(fileupload.FieldProvider, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.ExternalID(); ok {
		_spec.SetField(fileupload.FieldExternalID, field.TypeString, value)
	}
	if value, ok := fuuo.mutation.SizeBytes(); ok {
		_spec.SetField(fileupload.FieldSizeBytes, field.TypeInt64, value)
	}
	if value, ok := fuuo.mutation.AddedSizeBytes(); ok {
		_spec.AddField(fileupload.FieldSizeBytes, field.TypeInt64, value)
	}
	if value, ok := fuuo.mutation.FileSize(); ok {
		_spec.SetField(fileupload.FieldFileSize, field.TypeString, value)
	}
	if fuuo.mutation.FileSizeCleared() {
		_spec.ClearField(fileupload.FieldFileSize, field.TypeString)
	}
	_node = &FileUpload{config: fuuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fileupload.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fuuo.mutation.done = true
	return _node, nil
}
