package cache

import (
	"github.com/omkar273/codegeeky/internal/logger"
)

// Initialize initializes the cache system
func Initialize(log *logger.Logger) Cache {
	log.Info("Initializing cache system")

	return NewInMemoryCache()
}
