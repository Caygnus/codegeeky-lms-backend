package fileupload

import "context"

type Repository interface {
	Create(ctx context.Context, file *FileUpload) (*FileUpload, error)
	Get(ctx context.Context, id string) (*FileUpload, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}
