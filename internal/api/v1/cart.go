package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/service"
)

type CartHandler struct {
	cartService service.CartService
	logger      *logger.Logger
}

func NewCartHandler(cartService service.CartService, logger *logger.Logger) *CartHandler {
	return &CartHandler{cartService: cartService, logger: logger}
}

// @Summary Get a cart
// @Description Get a cart by ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} ierr.ErrorResponse
// @Failure 500 {object} ierr.ErrorResponse
// @Router /carts/{id} [get]
// @Security ApiKeyAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	id := c.Param("id")

	cart, err := h.cartService.GetCart(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get cart", "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cart)
}
