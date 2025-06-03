package fileupload

import (
	"context"
	"mime/multipart"

	"github.com/omkar273/codegeeky/internal/types"
)

type Provider interface {
	GetProvider() types.FileUploadProvider
	UploadFile(ctx context.Context, file *FileUpload, fileData *multipart.FileHeader) (*FileUpload, error)
	GetPresignedURL(ctx context.Context, externalID string) (string, error)
	DownloadFile(ctx context.Context, externalID string) ([]byte, error)
	Exists(ctx context.Context, externalID string) (bool, error)
	DeleteFile(ctx context.Context, externalID string) error
}

