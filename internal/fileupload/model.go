package fileupload

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type FileUpload struct {
	// basic file info
	ID        string              `json:"id,omitempty"`
	FileName  string              `json:"file_name,omitempty"`  // original file name
	FileType  string              `json:"file_type,omitempty"`  // e.g. image, video, document
	Extension types.FileExtension `json:"extension,omitempty"`  // e.g. .jpg, .png
	MimeType  string              `json:"mime_type,omitempty"`  // MIME type of the file (e.g. application/pdf, image/jpeg, video/mp4, etc.)
	PublicURL string              `json:"public_url,omitempty"` // public URL
	SecureURL *string             `json:"secure_url,omitempty"` // secure CDN URL or public URL

	// file upload info
	Provider   types.FileUploadProvider `json:"provider,omitempty"`
	ExternalID string                   `json:"external_id,omitempty"`

	// size info
	SizeBytes int64   `json:"size_bytes,omitempty"` // raw bytes count
	FileSize  *string `json:"file_size,omitempty"`  // human-readable e.g. "1.5 MB"

	types.BaseModel
}

func FromEnt(ent *ent.FileUpload) *FileUpload {
	return &FileUpload{
		ID:         ent.ID,
		FileName:   ent.FileName,
		FileType:   ent.FileType,
		Extension:  types.FileExtension(ent.Extension),
		MimeType:   ent.MimeType,
		PublicURL:  ent.PublicURL,
		SecureURL:  ent.SecureURL,
		Provider:   types.FileUploadProvider(ent.Provider),
		ExternalID: ent.ExternalID,
		SizeBytes:  ent.SizeBytes,
		FileSize:   ent.FileSize,
		BaseModel: types.BaseModel{
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
		},
	}
}

func FromEntList(ents []*ent.FileUpload) []*FileUpload {
	return lo.Map(ents, func(ent *ent.FileUpload, _ int) *FileUpload {
		return FromEnt(ent)
	})
}
