package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Get current user
// @Description Get the current user's information
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} dto.MeResponse
// @Failure 401 {object} ierr.Error
// @Failure 500 {object} ierr.Error
// @Router /user/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.userService.Me(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Update current user
// @Description Update the current user's information
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} dto.MeResponse
// @Failure 400 {object} ierr.Error
// @Failure 500 {object} ierr.Error
// @Router /user [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(ierr.WithError(err).
			WithHint("Please check the request payload").
			Mark(ierr.ErrValidation))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}
