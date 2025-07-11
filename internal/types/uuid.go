package types

import (
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

// GenerateUUID returns a k-sortable unique identifier
func GenerateUUID() string {
	return ulid.Make().String()
}

// GenerateUUIDWithPrefix returns a k-sortable unique identifier
// with a prefix ex inv_0ujsswThIGTUYm2K8FjOOfXtY1K
func GenerateUUIDWithPrefix(prefix string) string {
	if prefix == "" {
		return GenerateUUID()
	}
	return fmt.Sprintf("%s_%s", prefix, GenerateUUID())
}

func ValidateUUID(uuid string) bool {
	_, err := ulid.Parse(uuid)
	return err == nil
}

func ValidateUUIDWithPrefix(uuid string, prefix string) bool {
	if !strings.HasPrefix(uuid, prefix+"_") {
		return false
	}
	// Extract the ULID part after the prefix and underscore
	ulidPart := strings.TrimPrefix(uuid, prefix+"_")
	return ValidateUUID(ulidPart)
}

const (
	// Prefixes for all domains and entities
	UUID_PREFIX_USER                        = "user"
	UUID_PREFIX_INTERNSHIP                  = "internship"
	UUID_PREFIX_CATEGORY                    = "category"
	UUID_PREFIX_FILE_UPLOAD                 = "file"
	UUID_PREFIX_PAYMENT                     = "payment"
	UUID_PREFIX_DISCOUNT                    = "discount"
	UUID_PREFIX_DISCOUNT_APPLIED            = "da"
	UUID_PREFIX_INTERNSHIP_ENROLLMENT       = "ie"
	UUID_PREFIX_PAYMENT_ATTEMPT             = "pa"
	UUID_PREFIX_INTERNSHIP_BATCH            = "batch"
	UUID_PREFIX_CART                        = "cart"
	UUID_PREFIX_CART_LINE_ITEM              = "cli"
	UUID_PREFIX_INTERNSHIP_BATCH_ENROLLMENT = "enrollment"
)
