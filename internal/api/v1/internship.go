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

type InternshipHandler struct {
	internshipService service.InternshipService
	logger            *logger.Logger
}

func NewInternshipHandler(
	internshipService service.InternshipService,
	logger *logger.Logger,
) *InternshipHandler {
	return &InternshipHandler{
		internshipService: internshipService,
		logger:            logger,
	}
}

// @Summary Create a new internship
// @Description Create a new internship with the provided details
// @Tags Internship
// @Accept json
// @Produce json
// @Param internship body dto.CreateInternshipRequest true "Internship details"
// @Success 201 {object} dto.InternshipResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internships [post]
func (h *InternshipHandler) CreateInternship(c *gin.Context) {

	var req dto.CreateInternshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	internship, err := h.internshipService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorw("Failed to create internship", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, internship)
}

// @Summary Get an internship by ID
// @Description Get an internship by its unique identifier
// @Tags Internship
// @Accept json
// @Produce json
// @Param id path string true "Internship ID"
// @Success 200 {object} dto.InternshipResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internships/{id} [get]
func (h *InternshipHandler) GetInternship(c *gin.Context) {

	internshipID := c.Param("id")
	if internshipID == "" {
		c.Error(ierr.NewError("Internship ID is required").Mark(ierr.ErrValidation))
		return
	}

	internship, err := h.internshipService.GetByID(c.Request.Context(), internshipID)
	if err != nil {
		h.logger.Errorw("Failed to get internship", "error", err, "internship_id", internshipID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, internship)
}

// @Summary Update an internship
// @Description Update an internship by its unique identifier
// @Tags Internship
// @Accept json
// @Produce json
// @Param id path string true "Internship ID"
// @Param internship body dto.UpdateInternshipRequest true "Internship details"
// @Success 200 {object} dto.InternshipResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internships/{id} [put]
func (h *InternshipHandler) UpdateInternship(c *gin.Context) {

	internshipID := c.Param("id")
	if internshipID == "" {
		c.Error(ierr.NewError("Internship ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	var req dto.UpdateInternshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		return
	}

	internship, err := h.internshipService.Update(c.Request.Context(), internshipID, &req)
	if err != nil {
		h.logger.Errorw("Failed to update internship", "error", err, "internship_id", internshipID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, internship)
}

// @Summary Delete an internship
// @Description Delete an internship by its unique identifier
// @Tags Internship
// @Accept json
// @Produce json
// @Param id path string true "Internship ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internships/{id} [delete]
func (h *InternshipHandler) DeleteInternship(c *gin.Context) {

	internshipID := c.Param("id")
	if internshipID == "" {
		c.Error(ierr.NewError("Internship ID is required").Mark(ierr.ErrValidation))
		return
	}

	err := h.internshipService.Delete(c.Request.Context(), internshipID)
	if err != nil {
		h.logger.Errorw("Failed to delete internship", "error", err, "internship_id", internshipID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List internships
// @Description List internships with optional filtering
// @Tags Internship
// @Accept json
// @Produce json
// @Param filter query types.InternshipFilter true "Filter options"
// @Success 200 {object} dto.ListInternshipResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /internships [get]
func (h *InternshipHandler) ListInternships(c *gin.Context) {

	filter := types.NewInternshipFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.Error(err)
		return
	}

	if err := filter.Validate(); err != nil {
		c.Error(err)
		return
	}

	internships, err := h.internshipService.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Errorw("Failed to list internships", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, internships)
}
