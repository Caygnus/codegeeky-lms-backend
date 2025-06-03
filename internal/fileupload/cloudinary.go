package fileupload

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/omkar273/codegeeky/internal/config"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type cloudinaryProvider struct {
	baseProvider
	cfg        *config.Configuration
	logger     *logger.Logger
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryProvider(cfg *config.Configuration, logger *logger.Logger) (Provider, error) {

	cloudinaryClient, err := cloudinary.NewFromURL(cfg.Cloudinary.GetCloudinaryURL())
	if err != nil {
		return nil, ierr.WithError(err).
			WithHintf("Failed to create cloudinary client").
			Mark(ierr.ErrSystem)
	}

	return &cloudinaryProvider{
		baseProvider: baseProvider{
			provider: types.FileUploadProviderCloudinary,
		},
		cfg:        cfg,
		logger:     logger,
		cloudinary: cloudinaryClient,
	}, nil
}

func (p *cloudinaryProvider) GetProvider() types.FileUploadProvider {
	return p.baseProvider.GetProvider()
}

func (p *cloudinaryProvider) UploadFile(ctx context.Context, file *FileUpload, fileData *multipart.FileHeader) (*FileUpload, error) {
	// Open the uploaded file
	src, err := fileData.Open()
	if err != nil {
		return nil, ierr.WithError(err).
			WithHintf("Failed to open uploaded file").
			Mark(ierr.ErrSystem)
	}
	defer src.Close()

	// Extract file extension and set up upload parameters
	ext := strings.ToLower(filepath.Ext(fileData.Filename))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:] // Remove the leading dot
	}

	// Set up Cloudinary upload parameters
	useFilename := false
	uniqueFilename := false
	uploadParams := uploader.UploadParams{
		PublicID:       file.ID, // Use the FileUpload ID as the public ID
		ResourceType:   file.FileType,
		Folder:         "uploads",       // Optional: organize files in a folder
		UseFilename:    &useFilename,    // Use our custom public ID instead
		UniqueFilename: &uniqueFilename, // Don't add random characters
	}

	// Upload the file to Cloudinary
	uploadResult, err := p.cloudinary.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		p.logger.Error("Failed to upload file to Cloudinary", map[string]interface{}{
			"error":    err.Error(),
			"file_id":  file.ID,
			"filename": fileData.Filename,
		})
		return nil, ierr.WithError(err).
			WithHintf("Failed to upload file to Cloudinary").
			Mark(ierr.ErrSystem)
	}

	// Update the FileUpload struct with the upload result details
	file.ExternalID = uploadResult.PublicID
	file.PublicURL = uploadResult.URL
	file.SecureURL = lo.ToPtr(uploadResult.SecureURL)
	file.Provider = types.FileUploadProviderCloudinary

	// Update timestamps
	now := time.Now()
	file.UpdatedAt = now

	p.logger.Info("File uploaded successfully to Cloudinary", map[string]interface{}{
		"file_id":     file.ID,
		"external_id": file.ExternalID,
		"url":         file.PublicURL,
		"size_bytes":  file.SizeBytes,
	})

	return file, nil
}

func (p *cloudinaryProvider) DeleteFile(ctx context.Context, externalID string) error {
	_, err := p.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: externalID,
	})

	if err != nil {
		return ierr.WithError(err).
			WithHintf("Failed to delete file from Cloudinary").
			Mark(ierr.ErrSystem)
	}
	return nil
}
