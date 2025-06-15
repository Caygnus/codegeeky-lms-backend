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

type DiscountHandler struct {
	discountService service.DiscountService
	logger          *logger.Logger
}

func NewDiscountHandler(discountService service.DiscountService, logger *logger.Logger) *DiscountHandler {
	return &DiscountHandler{discountService: discountService, logger: logger}
}

// @Summary Create a new discount
// @Description Create a new discount with the provided details
// @Tags Discount
// @Accept json
// @Produce json
// @Param discount body dto.CreateDiscountRequest true "Discount details"
// @Success 201 {object} dto.DiscountResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts [post]
// @Security ApiKeyAuth
func (h *DiscountHandler) CreateDiscount(c *gin.Context) {

	var req dto.CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		c.Error(ierr.WithError(err).
			WithHint("Failed to bind request").
			Mark(ierr.ErrValidation))
		return
	}

	discount, err := h.discountService.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, discount)
}

// @Summary Get a discount by ID
// @Description Get a discount by its unique identifier
// @Tags Discount
// @Accept json
// @Produce json
// @Param id path string true "Discount ID"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts/{id} [get]
// @Security ApiKeyAuth
func (h *DiscountHandler) GetDiscount(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.Error(ierr.NewError("discount id is required").
			WithHint("Discount ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	discount, err := h.discountService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, discount)
}

// @Summary Update a discount by ID
// @Description Update a discount by its unique identifier
// @Tags Discount
// @Accept json
// @Produce json
// @Param id path string true "Discount ID"
// @Param discount body dto.UpdateDiscountRequest true "Discount details"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts/{id} [put]
// @Security ApiKeyAuth
func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.Error(ierr.NewError("discount id is required").
			WithHint("Discount ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	var req dto.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		c.Error(ierr.WithError(err).
			WithHint("Failed to bind request").
			Mark(ierr.ErrValidation))
		return
	}

	discount, err := h.discountService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, discount)
}

// @Summary Delete a discount by ID
// @Description Delete a discount by its unique identifier
// @Tags Discount
// @Accept json
// @Produce json
// @Param id path string true "Discount ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts/{id} [delete]
// @Security ApiKeyAuth
func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.Error(ierr.NewError("discount id is required").
			WithHint("Discount ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	if err := h.discountService.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List discounts
// @Description List discounts with optional filtering
// @Tags Discount
// @Accept json
// @Produce json
// @Param filter query types.DiscountFilter true "Filter options"
// @Success 200 {object} dto.ListDiscountResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts [get]
// @Security ApiKeyAuth
func (h *DiscountHandler) ListDiscounts(c *gin.Context) {
	filter := types.NewDiscountFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.Error(ierr.WithError(err).
			WithHint("Failed to bind query").
			Mark(ierr.ErrValidation))
		return
	}

	if err := filter.Validate(); err != nil {
		c.Error(err)
		return
	}

	discounts, err := h.discountService.List(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, discounts)
}

// @Summary Get a discount by code
// @Description Get a discount by its unique code
// @Tags Discount
// @Accept json
// @Produce json
// @Param code path string true "Discount code"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /discounts/code/{code} [get]
// @Security ApiKeyAuth
func (h *DiscountHandler) GetDiscountByCode(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.Error(ierr.NewError("discount code is required").
			WithHint("Discount code is required").
			Mark(ierr.ErrValidation))
		return
	}

	discount, err := h.discountService.GetByCode(c.Request.Context(), code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, discount)
}
