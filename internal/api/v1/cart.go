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

type CartHandler struct {
	cartService service.CartService
	logger      *logger.Logger
}

func NewCartHandler(cartService service.CartService, logger *logger.Logger) *CartHandler {
	return &CartHandler{cartService: cartService, logger: logger}
}

// @Summary Create a new cart
// @Description Create a new cart with the provided details
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart body dto.CreateCartRequest true "Cart details"
// @Success 201 {object} dto.CartResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts [post]
// @Security ApiKeyAuth
func (h *CartHandler) CreateCart(c *gin.Context) {
	var req dto.CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cart, err := h.cartService.CreateCart(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorw("Failed to create cart", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, cart)
}

// @Summary Get a cart
// @Description Get a cart by ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id} [get]
// @Security ApiKeyAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(ierr.NewError("Cart ID is required").Mark(ierr.ErrValidation))
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get cart", "error", err, "cart_id", id)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cart)
}

// @Summary Update a cart
// @Description Update a cart by its unique identifier
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param cart body dto.UpdateCartRequest true "Cart details"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id} [put]
// @Security ApiKeyAuth
func (h *CartHandler) UpdateCart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(ierr.NewError("Cart ID is required").Mark(ierr.ErrValidation))
		return
	}

	var req dto.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cart, err := h.cartService.UpdateCart(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Errorw("Failed to update cart", "error", err, "cart_id", id)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cart)
}

// @Summary Delete a cart
// @Description Delete a cart by its unique identifier
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id} [delete]
// @Security ApiKeyAuth
func (h *CartHandler) DeleteCart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(ierr.NewError("Cart ID is required").Mark(ierr.ErrValidation))
		return
	}

	err := h.cartService.DeleteCart(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to delete cart", "error", err, "cart_id", id)
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List carts
// @Description List carts with optional filtering
// @Tags Cart
// @Accept json
// @Produce json
// @Param filter query types.CartFilter true "Filter options"
// @Success 200 {object} dto.ListCartResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts [get]
// @Security ApiKeyAuth
func (h *CartHandler) ListCarts(c *gin.Context) {
	filter := types.NewCartFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.Error(err)
		return
	}

	if err := filter.Validate(); err != nil {
		c.Error(err)
		return
	}

	carts, err := h.cartService.ListCarts(c.Request.Context(), filter)
	if err != nil {
		h.logger.Errorw("Failed to list carts", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, carts)
}

// @Summary Get cart line items
// @Description Get all line items for a specific cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {array} cart.CartLineItem
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id}/line-items [get]
// @Security ApiKeyAuth
func (h *CartHandler) GetCartLineItems(c *gin.Context) {
	cartID := c.Param("id")
	if cartID == "" {
		c.Error(ierr.NewError("Cart ID is required").Mark(ierr.ErrValidation))
		return
	}

	lineItems, err := h.cartService.GetCartLineItems(c.Request.Context(), cartID)
	if err != nil {
		h.logger.Errorw("Failed to get cart line items", "error", err, "cart_id", cartID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lineItems)
}

// @Summary Add line item to cart
// @Description Add a new line item to a cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param line_item body dto.CreateCartLineItemRequest true "Line item details"
// @Success 201 {object} cart.CartLineItem
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id}/line-items [post]
// @Security ApiKeyAuth
func (h *CartHandler) AddLineItem(c *gin.Context) {
	cartID := c.Param("id")
	if cartID == "" {
		c.Error(ierr.NewError("Cart ID is required").Mark(ierr.ErrValidation))
		return
	}

	var req dto.CreateCartLineItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	// Set the cart ID from the URL parameter
	req.CartID = cartID

	lineItem, err := h.cartService.AddLineItem(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorw("Failed to add line item", "error", err, "cart_id", cartID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, lineItem)
}

// @Summary Remove line item from cart
// @Description Remove a line item from a cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param line_item_id path string true "Line Item ID"
// @Success 204
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id}/line-items/{line_item_id} [delete]
// @Security ApiKeyAuth
func (h *CartHandler) RemoveLineItem(c *gin.Context) {
	lineItemID := c.Param("line_item_id")
	if lineItemID == "" {
		c.Error(ierr.NewError("Line item ID is required").Mark(ierr.ErrValidation))
		return
	}

	err := h.cartService.RemoveLineItem(c.Request.Context(), lineItemID)
	if err != nil {
		h.logger.Errorw("Failed to remove line item", "error", err, "line_item_id", lineItemID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary Get line item
// @Description Get a specific line item by its ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param line_item_id path string true "Line Item ID"
// @Success 200 {object} cart.CartLineItem
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id}/line-items/{line_item_id} [get]
// @Security ApiKeyAuth
func (h *CartHandler) GetLineItem(c *gin.Context) {
	lineItemID := c.Param("line_item_id")
	if lineItemID == "" {
		c.Error(ierr.NewError("Line item ID is required").Mark(ierr.ErrValidation))
		return
	}

	lineItem, err := h.cartService.GetLineItem(c.Request.Context(), lineItemID)
	if err != nil {
		h.logger.Errorw("Failed to get line item", "error", err, "line_item_id", lineItemID)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lineItem)
}

// @Summary Get user's default cart
// @Description Get the default cart for the authenticated user
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {object} dto.CartResponse
// @Failure 404 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/default [get]
// @Security ApiKeyAuth
func (h *CartHandler) GetDefaultCart(c *gin.Context) {
	cart, err := h.cartService.GetUserDefaultCart(c.Request.Context())
	if err != nil {
		h.logger.Errorw("Failed to get default cart", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cart)
}
