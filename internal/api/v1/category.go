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

type CategoryHandler struct {
	categoryService service.CategoryService
	logger          *logger.Logger
}

func NewCategoryHandler(categoryService service.CategoryService, logger *logger.Logger) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, logger: logger}
}

// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags Category
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorw("Failed to create category", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Get a category by ID
// @Description Get a category by its unique identifier
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {

	categoryID := c.Param("id")
	if categoryID == "" {
		c.Error(ierr.NewError("Category ID is required").
			WithHint("Category ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), categoryID)
	if err != nil {
		h.logger.Errorw("Failed to get category", "error", err, "category_id", categoryID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Update a category
// @Description Update a category by its unique identifier
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body dto.UpdateCategoryRequest true "Category details"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {

	categoryID := c.Param("id")
	if categoryID == "" {
		c.Error(ierr.NewError("Category ID is required").
			WithHint("Category ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
}

// @Summary Delete a category
// @Description Delete a category by its unique identifier
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID := c.Param("id")
	if categoryID == "" {
		c.Error(ierr.NewError("Category ID is required").
			WithHint("Category ID is required").
			Mark(ierr.ErrValidation))
		return
	}

	err := h.categoryService.Delete(c.Request.Context(), categoryID)
	if err != nil {
		h.logger.Errorw("Failed to delete category", "error", err, "category_id", categoryID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List categories
// @Description List categories with optional filtering
// @Tags Category
// @Accept json
// @Produce json
// @Param filter query types.CategoryFilter true "Filter options"
// @Success 200 {object} dto.ListCategoryResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {

	filter := types.NewCategoryFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.Error(err)
		return
	}

	if err := filter.Validate(); err != nil {
		c.Error(err)
		return
	}

	categories, err := h.categoryService.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Errorw("Failed to list categories", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, categories)
}
