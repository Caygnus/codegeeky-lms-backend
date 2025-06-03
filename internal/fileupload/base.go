package fileupload

import (
	"context"
	"errors"
	"mime/multipart"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type baseProvider struct {
	provider types.FileUploadProvider
}

func NewBaseProvider() (Provider, error) {
	return &baseProvider{
		provider: types.FileUploadProviderS3,
	}, nil
}

func (p *baseProvider) GetProvider() types.FileUploadProvider {
	return p.provider
}

func (p *baseProvider) UploadFile(ctx context.Context, file *FileUpload, fileData *multipart.FileHeader) (*FileUpload, error) {
	return nil, ierr.WithError(errors.New("not implemented")).
		WithHintf("UploadFile is not implemented for provider: %s", p.GetProvider()).
		Mark(ierr.ErrSystem)
}

func (p *baseProvider) GetPresignedURL(ctx context.Context, externalID string) (string, error) {
	return "", ierr.WithError(errors.New("not implemented")).
		WithHintf("GetPresignedURL is not implemented for provider: %s", p.GetProvider()).
		Mark(ierr.ErrSystem)
}

func (p *baseProvider) DownloadFile(ctx context.Context, externalID string) ([]byte, error) {
	return nil, ierr.WithError(errors.New("not implemented")).
		WithHintf("DownloadFile is not implemented for provider: %s", p.GetProvider()).
		Mark(ierr.ErrSystem)
}

func (p *baseProvider) Exists(ctx context.Context, externalID string) (bool, error) {
	return false, ierr.WithError(errors.New("not implemented")).
		WithHintf("Exists is not implemented for provider: %s", p.GetProvider()).
		Mark(ierr.ErrSystem)
}

func (p *baseProvider) DeleteFile(ctx context.Context, externalID string) error {
	return ierr.WithError(errors.New("not implemented")).
		WithHintf("DeleteFile is not implemented for provider: %s", p.GetProvider()).
		Mark(ierr.ErrSystem)
}
