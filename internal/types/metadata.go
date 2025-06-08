package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	ierr "github.com/omkar273/codegeeky/internal/errors"
)

// Metadata represents a JSONB field for storing key-value pairs
type Metadata map[string]string

// convertEntMetadata safely converts Ent JSONMap to Metadata
func MetadataFromEnt(src map[string]string) Metadata {
	meta := make(Metadata, len(src))
	for k, v := range src {
		meta[k] = fmt.Sprintf("%v", v)
	}
	return meta
}

// Scan implements the sql.Scanner interface for Metadata
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = make(Metadata)
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return ierr.NewError("failed to scan metadata").
			WithHint(fmt.Sprintf("Expected []byte or string but got %T", value)).
			Mark(ierr.ErrValidation)
	}

	result := make(Metadata)
	if err := json.Unmarshal(data, &result); err != nil {
		return ierr.NewError("failed to parse metadata JSON").
			WithHint("Ensure metadata is valid JSON").
			WithReportableDetails(map[string]any{
				"error": err.Error(),
			}).
			Mark(ierr.ErrValidation)
	}

	*m = result
	return nil
}

// Value implements the driver.Valuer interface for Metadata
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(Metadata{})
	}
	return json.Marshal(m)
}
