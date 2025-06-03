package types

type FileType string

const (
	FileTypePDF         FileType = "pdf"
	FileTypeImage       FileType = "image"
	FileTypeVideo       FileType = "video"
	FileTypeDocument    FileType = "document"
	FileTypeResume      FileType = "resume"
	FileTypeCoverLetter FileType = "cover_letter"
	FileTypeTranscript  FileType = "transcript"
	FileTypeCertificate FileType = "certificate"
	FileTypeOther       FileType = "other" // for other file types
)

type FileUploadProvider string

const (
	FileUploadProviderCloudinary FileUploadProvider = "cloudinary"
	FileUploadProviderS3         FileUploadProvider = "s3"
)
