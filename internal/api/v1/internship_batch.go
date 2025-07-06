package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/service"
	"github.com/omkar273/codegeeky/internal/types"
)

type InternshipBatchHandler struct {
	internshipBatchService service.InternshipBatchService
	logger                 *logger.Logger
}

func NewInternshipBatchHandler(
	internshipBatchService service.InternshipBatchService,
	logger *logger.Logger,
) *InternshipBatchHandler {
	return &InternshipBatchHandler{
		internshipBatchService: internshipBatchService,
		logger:                 logger,
	}
}

// @Summary Create a new internship batch
// @Description Create a new internship batch with the provided details
// @Tags InternshipBatch
// @Accept json
// @Produce json
// @Param batch body dto.CreateInternshipBatchRequest true "Internship batch details"
// @Success 201 {object} dto.InternshipBatchResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internshipbatches [post]
func (h *InternshipBatchHandler) CreateInternshipBatch(c *gin.Context) {
	var req dto.CreateInternshipBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	batch, err := h.internshipBatchService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorw("Failed to create internship batch", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, batch)
}

// @Summary Get an internship batch by ID
// @Description Get an internship batch by its unique identifier
// @Tags InternshipBatch
// @Accept json
// @Produce json
// @Param id path string true "Internship Batch ID"
// @Success 200 {object} dto.InternshipBatchResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internshipbatches/{id} [get]
func (h *InternshipBatchHandler) GetInternshipBatch(c *gin.Context) {
	batchID := c.Param("id")
	if batchID == "" {
		c.Error(ierr.NewError("Internship batch ID is required").
			WithHint("Internship batch ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	batch, err := h.internshipBatchService.Get(c.Request.Context(), batchID)
	if err != nil {
		h.logger.Errorw("Failed to get internship batch", "error", err, "batch_id", batchID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, batch)
}

// @Summary Update an internship batch
// @Description Update an internship batch by its unique identifier
// @Tags InternshipBatch
// @Accept json
// @Produce json
// @Param id path string true "Internship Batch ID"
// @Param batch body dto.UpdateInternshipBatchRequest true "Internship batch details"
// @Success 200 {object} dto.InternshipBatchResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internshipbatches/{id} [put]
func (h *InternshipBatchHandler) UpdateInternshipBatch(c *gin.Context) {
	batchID := c.Param("id")
	if batchID == "" {
		c.Error(ierr.NewError("Internship batch ID is required").
			WithHint("Internship batch ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	var req dto.UpdateInternshipBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		return
	}

	batch, err := h.internshipBatchService.Update(c.Request.Context(), batchID, &req)
	if err != nil {
		h.logger.Errorw("Failed to update internship batch", "error", err, "batch_id", batchID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, batch)
}

// @Summary Delete an internship batch
// @Description Delete an internship batch by its unique identifier
// @Tags InternshipBatch
// @Accept json
// @Produce json
// @Param id path string true "Internship Batch ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internshipbatches/{id} [delete]
func (h *InternshipBatchHandler) DeleteInternshipBatch(c *gin.Context) {
	batchID := c.Param("id")
	if batchID == "" {
		c.Error(ierr.NewError("Internship batch ID is required").
			WithHint("Internship batch ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	err := h.internshipBatchService.Delete(c.Request.Context(), batchID)
	if err != nil {
		h.logger.Errorw("Failed to delete internship batch", "error", err, "batch_id", batchID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List internship batches
// @Description List internship batches with optional filtering
// @Tags InternshipBatch
// @Accept json
// @Produce json
// @Param filter query types.InternshipBatchFilter true "Filter options"
// @Success 200 {object} dto.ListInternshipBatchResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internshipbatches [get]
func (h *InternshipBatchHandler) ListInternshipBatches(c *gin.Context) {
	filter := types.NewInternshipBatchFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.Error(err)
		return
	}

	if err := filter.Validate(); err != nil {
		c.Error(err)
		return
	}

	batches, err := h.internshipBatchService.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Errorw("Failed to list internship batches", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, batches)
}
