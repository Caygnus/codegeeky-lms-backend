package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type FileExtractOptions struct {
	FieldName   string                // e.g. "file"
	MaxSize     int64                 // in bytes
	AllowedExts []types.FileExtension // e.g. []types.FileExtension{types.FileExtensionJPG, types.FileExtensionPNG, types.FileExtensionPDF}
}

type FileExtractResult struct {
	FileName  string              `json:"file_name"`
	FileType  string              `json:"file_type"`
	Extension types.FileExtension `json:"extension"`
	MimeType  string              `json:"mime_type"`
	SizeBytes int64               `json:"size_bytes"`
	FileSize  string              `json:"file_size"`

	// file upload info
	FileBuffer *multipart.File
}

func NewFileExtractOptions(fieldName string, maxSize int64, allowedExts []types.FileExtension) *FileExtractOptions {
	return &FileExtractOptions{
		FieldName:   fieldName,
		MaxSize:     maxSize,
		AllowedExts: allowedExts,
	}
}

// ExtractFile extracts a file from the request and returns a fileupload.FileUpload and a multipart.File
// It sets the file upload fields:
// - FileName: the original file name
// - FileType: the type of the file (image, video, document, other)
// - Extension: the extension of the file
// - MimeType: the MIME type of the file
// - SizeBytes: the size of the file in bytes
// - FileSize: the size of the file in a human-readable format
//
// It returns the file upload object and the file as a multipart.File
// The file upload object is the file upload object that will be used to upload the file to the provider
// The file is the file that was extracted from the request
func ExtractFile(
	c *gin.Context,
	opts FileExtractOptions,
) (*FileExtractResult, error) {
	if opts.FieldName == "" {
		return nil, ierr.NewError("missing file field name").
			Mark(ierr.ErrValidation)
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, opts.MaxSize)

	if err := c.Request.ParseMultipartForm(opts.MaxSize); err != nil {
		return nil, ierr.WithError(err).
			WithHintf("Failed to parse multipart form, likely exceeded max size of %.2f MB", float64(opts.MaxSize)/float64(types.MB)).
			Mark(ierr.ErrFileTooLarge)
	}

	fileHeader, err := c.FormFile(opts.FieldName)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("File not found in request").
			Mark(ierr.ErrNotFound)
	}

	ext := types.FileExtension(strings.ToLower(filepath.Ext(fileHeader.Filename)))
	normalizedExts := slices.Clone(opts.AllowedExts)
	for i, e := range normalizedExts {
		normalizedExts[i] = types.FileExtension(strings.ToLower(string(e)))
	}
	if !slices.Contains(normalizedExts, ext) {
		return nil, ierr.NewError("invalid file extension").
			WithHintf("Extension %q is not allowed. Allowed: %v", ext, opts.AllowedExts).
			Mark(ierr.ErrInvalidExtension)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to open uploaded file").
			Mark(ierr.ErrSystem)
	}

	return &FileExtractResult{
		FileName:   fileHeader.Filename,
		FileType:   detectFileType(ext),
		Extension:  ext,
		MimeType:   fileHeader.Header.Get("Content-Type"),
		SizeBytes:  fileHeader.Size,
		FileSize:   humanFileSize(fileHeader.Size),
		FileBuffer: &file,
	}, nil
}

func humanFileSize(size int64) string {
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	div := float64(unit)
	s := float64(size)

	for _, u := range units {
		if s < div*1000 {
			return fmt.Sprintf("%.2f %s", s/div, u)
		}
		s /= 1000
		div *= 1000
	}

	// For exabytes and beyond
	return fmt.Sprintf("%.2f ZB", s/1000)
}

func detectFileType(ext types.FileExtension) string {
	switch ext {
	case types.FileExtensionJPG, types.FileExtensionJPEG, types.FileExtensionPNG, types.FileExtensionGIF, types.FileExtensionWEBP, types.FileExtensionSVG:
		return "image"
	case types.FileExtensionPDF, types.FileExtensionDOC, types.FileExtensionPPT, types.FileExtensionPPTX, types.FileExtensionXLS, types.FileExtensionXLSX:
		return "document"
	case types.FileExtensionMP4, types.FileExtensionMOV, types.FileExtensionAVI, types.FileExtensionMKV:
		return "video"
	default:
		return "other"
	}
}
