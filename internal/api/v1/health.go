package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/logger"
)

type HealthHandler struct {
	logger *logger.Logger
}

func NewHealthHandler(logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// @Summary Health check
// @Description Health check
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
func (h *HealthHandler) Health(c *gin.Context) {
	request := map[string]interface{}{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	// log the incoming request
	h.logger.Infof("incoming request: %s %s | body: %s", c.Request.Method, c.Request.URL.Path, request)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
