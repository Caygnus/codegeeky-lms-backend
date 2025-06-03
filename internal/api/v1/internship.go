package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/auth"
	domainAuth "github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/rest/middleware"
	"github.com/omkar273/codegeeky/internal/service"
	"github.com/omkar273/codegeeky/internal/types"
)

type InternshipHandler struct {
	internshipService service.InternshipService
	authzService      auth.AuthorizationService
	logger            *logger.Logger
}

func NewInternshipHandler(
	internshipService service.InternshipService,
	authzService auth.AuthorizationService,
	logger *logger.Logger,
) *InternshipHandler {
	return &InternshipHandler{
		internshipService: internshipService,
		authzService:      authzService,
		logger:            logger,
	}
}

// CreateInternship creates a new internship
// This endpoint requires instructor or admin role (RBAC)
func (h *InternshipHandler) CreateInternship(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse request
	var req dto.CreateInternshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create internship using service (which handles authorization)
	internship, err := h.internshipService.Create(c.Request.Context(), &req, authContext.UserID, authContext.Role)
	if err != nil {
		h.logger.Errorw("Failed to create internship", "error", err, "user_id", authContext.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create internship"})
		return
	}

	// Return response
	response := &dto.InternshipResponse{}
	c.JSON(http.StatusCreated, response.FromDomain(internship))
}

// GetInternship retrieves an internship by ID
// This endpoint uses ABAC for fine-grained access control
func (h *InternshipHandler) GetInternship(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get internship ID from URL
	internshipID := c.Param("id")
	if internshipID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internship ID is required"})
		return
	}

	// Get internship using service (which handles authorization)
	internship, err := h.internshipService.GetByID(c.Request.Context(), internshipID, authContext)
	if err != nil {
		h.logger.Errorw("Failed to get internship", "error", err, "user_id", authContext.UserID, "internship_id", internshipID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Internship not found"})
		return
	}

	// Return response
	response := &dto.InternshipResponse{}
	c.JSON(http.StatusOK, response.FromDomain(internship))
}

// UpdateInternship updates an existing internship
// This endpoint uses ABAC to check ownership and permissions
func (h *InternshipHandler) UpdateInternship(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get internship ID from URL
	internshipID := c.Param("id")
	if internshipID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internship ID is required"})
		return
	}

	// Parse request
	var req dto.UpdateInternshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update internship using service (which handles authorization)
	internship, err := h.internshipService.Update(c.Request.Context(), internshipID, &req, authContext)
	if err != nil {
		h.logger.Errorw("Failed to update internship", "error", err, "user_id", authContext.UserID, "internship_id", internshipID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update internship"})
		return
	}

	// Return response
	response := &dto.InternshipResponse{}
	c.JSON(http.StatusOK, response.FromDomain(internship))
}

// DeleteInternship deletes an internship
// This endpoint uses ABAC to check ownership and permissions
func (h *InternshipHandler) DeleteInternship(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get internship ID from URL
	internshipID := c.Param("id")
	if internshipID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internship ID is required"})
		return
	}

	// Delete internship using service (which handles authorization)
	err := h.internshipService.Delete(c.Request.Context(), internshipID, authContext)
	if err != nil {
		h.logger.Errorw("Failed to delete internship", "error", err, "user_id", authContext.UserID, "internship_id", internshipID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete internship"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListInternships lists internships with role-based filtering
// This endpoint applies different filters based on user role
func (h *InternshipHandler) ListInternships(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse query parameters into filter
	filter := types.NewInternshipFilter()
	if err := c.ShouldBindQuery(filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate filter
	if err := filter.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// List internships using service (which handles role-based filtering)
	internships, err := h.internshipService.List(c.Request.Context(), filter, authContext)
	if err != nil {
		h.logger.Errorw("Failed to list internships", "error", err, "user_id", authContext.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list internships"})
		return
	}

	// Convert to response format
	responses := make([]*dto.InternshipResponse, len(internships))
	for i, internship := range internships {
		response := &dto.InternshipResponse{}
		responses[i] = response.FromDomain(internship)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"count": len(responses),
	})
}

// AccessInternshipContent demonstrates ABAC for content access
// This would be used for accessing lectures, assignments, etc.
func (h *InternshipHandler) AccessInternshipContent(c *gin.Context) {
	// Get auth context from middleware
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get internship and content IDs from URL
	internshipID := c.Param("internship_id")
	contentID := c.Param("content_id")
	contentType := c.Query("type") // lecture, assignment, resource

	if internshipID == "" || contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internship ID and Content ID are required"})
		return
	}

	// Create authorization request for content access
	authRequest := &domainAuth.AccessRequest{
		Subject: authContext,
		Resource: &domainAuth.Resource{
			Type: "content",
			ID:   contentID,
			Attributes: map[string]interface{}{
				"internship_id":     internshipID,
				"content_type":      contentType,
				"required_progress": 0.0, // This would come from database
			},
		},
		Action: domainAuth.PermissionViewLectures, // This would vary based on content type
	}

	// Check authorization using ABAC
	allowed, err := h.authzService.IsAuthorized(c.Request.Context(), authRequest)
	if err != nil {
		h.logger.Errorw("Authorization check failed for content access", "error", err, "user_id", authContext.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization check failed"})
		return
	}

	if !allowed {
		h.logger.Warnw("User not authorized to access content", "user_id", authContext.UserID, "internship_id", internshipID, "content_id", contentID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this content"})
		return
	}

	// If authorized, return the content (this would fetch from your content service)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Content access granted",
		"internship_id": internshipID,
		"content_id":    contentID,
		"content_type":  contentType,
	})
}
