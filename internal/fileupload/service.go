package fileupload

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Service interface {
	UploadFile(ctx context.Context, file *FileUpload) error
	GetPresignedUrl(ctx context.Context, id string, fileType types.FileType) (string, error)
	GetFile(ctx context.Context, id string, fileType types.FileType) ([]byte, error)
	Exists(ctx context.Context, id string, fileType types.FileType) (bool, error)
	DeleteFile(ctx context.Context, id string, fileType types.FileType) error
}
