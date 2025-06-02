package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Signup
// @Description Signup
// @Tags Auth
// @Accept json
// @Produce json
// @Param signupRequest body dto.SignupRequest true "Signup request"
// @Success 201 {object} dto.SignupResponse
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {

	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(ierr.WithError(err).
			WithHint("Please check the request payload").
			Mark(ierr.ErrValidation))
		return
	}

	resp, err := h.authService.Signup(c.Request.Context(), &req)
	if err != nil {
		c.Error(ierr.WithError(err).
			WithHint("Failed to signup").
			Mark(ierr.ErrDatabase))
		return
	}

	c.JSON(http.StatusCreated, resp)

}
