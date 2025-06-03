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

// File size helpers
const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

type FileExtension string

const (
	// image
	FileExtensionJPG  FileExtension = ".jpg"  // image/jpeg, image/pjpeg
	FileExtensionJPEG FileExtension = ".jpeg" // image/jpeg, image/pjpeg
	FileExtensionPNG  FileExtension = ".png"  // image/png, image/x-png
	FileExtensionGIF  FileExtension = ".gif"  // image/gif
	FileExtensionWEBP FileExtension = ".webp" // image/webp
	FileExtensionSVG  FileExtension = ".svg"  // image/svg+xml

	// document
	FileExtensionPDF  FileExtension = ".pdf"  // application/pdf
	FileExtensionDOC  FileExtension = ".doc"  // application/msword
	FileExtensionPPT  FileExtension = ".ppt"  // application/vnd.ms-powerpoint
	FileExtensionPPTX FileExtension = ".pptx" // application/vnd.openxmlformats-officedocument.presentationml.presentation
	FileExtensionXLS  FileExtension = ".xls"  // application/vnd.ms-excel
	FileExtensionXLSX FileExtension = ".xlsx" // application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
	FileExtensionTXT  FileExtension = ".txt"  // text/plain

	// video
	FileExtensionMP4 FileExtension = ".mp4" // video/mp4
	FileExtensionMOV FileExtension = ".mov" // video/quicktime
	FileExtensionAVI FileExtension = ".avi" // video/x-msvideo
	FileExtensionMKV FileExtension = ".mkv" // video/x-matroska

	// other
	FileExtensionZIP FileExtension = ".zip" // application/zip
)

type FileUploadFilter struct {
	*QueryFilter
	*TimeRangeFilter

	// These fields are used to filter file uploads by external id (id of the file upload in the provider)
	ExternalIDs []string `json:"external_ids,omitempty" form:"external_ids" validate:"omitempty"`

	// These fields are used to filter file uploads by file id (internal id of the file upload)
	FileIDs []string `json:"file_ids,omitempty" form:"file_ids" validate:"omitempty"`
}

func (f *FileUploadFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	return nil
}

func NewFileUploadFilter() *FileUploadFilter {
	return &FileUploadFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitFileUploadFilter() *FileUploadFilter {
	return &FileUploadFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *FileUploadFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *FileUploadFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *FileUploadFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *FileUploadFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *FileUploadFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *FileUploadFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *FileUploadFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
