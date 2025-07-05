package testutil

import (
	"context"
	"time"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/fileupload"
	"github.com/omkar273/codegeeky/internal/types"
)

// InMemoryFileUploadStore implements fileupload.Repository
type InMemoryFileUploadStore struct {
	*InMemoryStore[*fileupload.FileUpload]
}

// NewInMemoryFileUploadStore creates a new in-memory file upload store
func NewInMemoryFileUploadStore() *InMemoryFileUploadStore {
	return &InMemoryFileUploadStore{
		InMemoryStore: NewInMemoryStore[*fileupload.FileUpload](),
	}
}

func (s *InMemoryFileUploadStore) Create(ctx context.Context, f *fileupload.FileUpload) (*fileupload.FileUpload, error) {
	if f == nil {
		return nil, ierr.NewError("file upload cannot be nil").
			WithHint("File upload data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	if f.UpdatedAt.IsZero() {
		f.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, f.ID, f)
	if err != nil {
		if err.Error() == "item already exists" {
			return nil, ierr.WithError(err).
				WithHint("A file upload with this ID already exists").
				WithReportableDetails(map[string]any{
					"file_id":     f.ID,
					"external_id": f.ExternalID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to create file upload").
			Mark(ierr.ErrDatabase)
	}
	return f, nil
}

func (s *InMemoryFileUploadStore) Get(ctx context.Context, id string) (*fileupload.FileUpload, error) {
	file, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("File upload with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"file_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get file upload with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return file, nil
}

func (s *InMemoryFileUploadStore) Delete(ctx context.Context, id string) error {
	// Get the file first
	f, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	f.Status = types.StatusDeleted
	f.UpdatedAt = time.Now().UTC()

	return s.InMemoryStore.Update(ctx, f.ID, f)
}

func (s *InMemoryFileUploadStore) Exists(ctx context.Context, id string) (bool, error) {
	_, err := s.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Clear clears the file upload store
func (s *InMemoryFileUploadStore) Clear() {
	s.InMemoryStore.Clear()
}
