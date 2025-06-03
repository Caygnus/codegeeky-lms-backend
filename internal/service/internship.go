package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/auth"
	domainAuth "github.com/omkar273/codegeeky/internal/domain/auth"
	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type InternshipService interface {
	Create(ctx context.Context, req *dto.CreateInternshipRequest, userID string, userRole types.UserRole) (*domainInternship.Internship, error)
	GetByID(ctx context.Context, id string, authContext *domainAuth.AuthContext) (*domainInternship.Internship, error)
	Update(ctx context.Context, id string, req *dto.UpdateInternshipRequest, authContext *domainAuth.AuthContext) (*domainInternship.Internship, error)
	Delete(ctx context.Context, id string, authContext *domainAuth.AuthContext) error
	List(ctx context.Context, filter *types.InternshipFilter, authContext *domainAuth.AuthContext) ([]*domainInternship.Internship, error)
}

type internshipService struct {
	internshipRepo domainInternship.InternshipRepository
	authzService   auth.AuthorizationService
	logger         *logger.Logger
}

func NewInternshipService(
	internshipRepo domainInternship.InternshipRepository,
	authzService auth.AuthorizationService,
	logger *logger.Logger,
) InternshipService {
	return &internshipService{
		internshipRepo: internshipRepo,
		authzService:   authzService,
		logger:         logger,
	}
}

func (s *internshipService) Create(ctx context.Context, req *dto.CreateInternshipRequest, userID string, userRole types.UserRole) (*domainInternship.Internship, error) {
	// Create authorization request
	authRequest := &domainAuth.AccessRequest{
		Subject: &domainAuth.AuthContext{
			UserID: userID,
			Role:   userRole,
		},
		Resource: &domainAuth.Resource{
			Type: "internship",
		},
		Action: domainAuth.PermissionCreateInternship,
	}

	// Check authorization
	allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
	if err != nil {
		s.logger.Errorw("Authorization check failed", "error", err, "user_id", userID)
		return nil, err
	}

	if !allowed {
		s.logger.Warnw("User not authorized to create internship", "user_id", userID, "role", userRole)
		return nil, ierr.ErrPermissionDenied
	}

	// Create the internship (implement your business logic here)
	internship := &domainInternship.Internship{
		// Map from DTO to domain model
		Title:              req.Title,
		Description:        req.Description,
		Skills:             req.Skills,
		Level:              req.Level,
		Mode:               req.Mode,
		DurationInWeeks:    req.DurationInWeeks,
		LearningOutcomes:   req.LearningOutcomes,
		Prerequisites:      req.Prerequisites,
		Benefits:           req.Benefits,
		Currency:           req.Currency,
		Price:              req.Price,
		FlatDiscount:       lo.FromPtr(req.FlatDiscount),
		PercentageDiscount: lo.FromPtr(req.PercentageDiscount),
		BaseModel: types.BaseModel{
			Status:    types.StatusPublished,
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}

	// Save to repository
	err = s.internshipRepo.Create(ctx, internship)
	if err != nil {
		s.logger.Errorw("Failed to create internship", "error", err, "user_id", userID)
		return nil, err
	}

	s.logger.Infow("Internship created successfully", "internship_id", internship.ID, "user_id", userID)
	return internship, nil
}

func (s *internshipService) GetByID(ctx context.Context, id string, authContext *domainAuth.AuthContext) (*domainInternship.Internship, error) {
	// Get internship first to check ownership and other attributes
	internship, err := s.internshipRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if internship == nil {
		return nil, ierr.ErrNotFound
	}

	// Create authorization request
	authRequest := &domainAuth.AccessRequest{
		Subject: authContext,
		Resource: &domainAuth.Resource{
			Type: "internship",
			ID:   id,
			Attributes: map[string]interface{}{
				"created_by": internship.CreatedBy,
				"status":     internship.Status,
			},
		},
		Action: domainAuth.PermissionViewInternship,
	}

	// Check authorization
	allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
	if err != nil {
		s.logger.Errorw("Authorization check failed", "error", err, "user_id", authContext.UserID)
		return nil, err
	}

	if !allowed {
		s.logger.Warnw("User not authorized to view internship", "user_id", authContext.UserID, "internship_id", id)
		return nil, ierr.ErrPermissionDenied
	}

	return internship, nil
}

func (s *internshipService) Update(ctx context.Context, id string, req *dto.UpdateInternshipRequest, authContext *domainAuth.AuthContext) (*domainInternship.Internship, error) {
	// Get existing internship
	existingInternship, err := s.internshipRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingInternship == nil {
		return nil, ierr.ErrNotFound
	}

	// Create authorization request
	authRequest := &domainAuth.AccessRequest{
		Subject: authContext,
		Resource: &domainAuth.Resource{
			Type: "internship",
			ID:   id,
			Attributes: map[string]interface{}{
				"created_by": existingInternship.CreatedBy,
				"status":     existingInternship.Status,
			},
		},
		Action: domainAuth.PermissionUpdateInternship,
	}

	// Check authorization
	allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
	if err != nil {
		s.logger.Errorw("Authorization check failed", "error", err, "user_id", authContext.UserID)
		return nil, err
	}

	if !allowed {
		s.logger.Warnw("User not authorized to update internship", "user_id", authContext.UserID, "internship_id", id)
		return nil, ierr.ErrPermissionDenied
	}

	// Update the internship (implement your update logic here)
	if req.Title != nil {
		existingInternship.Title = *req.Title
	}
	if req.Description != nil {
		existingInternship.Description = *req.Description
	}
	if req.Level != nil {
		existingInternship.Level = *req.Level
	}
	if req.Mode != nil {
		existingInternship.Mode = *req.Mode
	}
	if req.DurationInWeeks != nil {
		existingInternship.DurationInWeeks = *req.DurationInWeeks
	}
	if req.Currency != nil {
		existingInternship.Currency = *req.Currency
	}
	if req.Price != nil {
		existingInternship.Price = *req.Price
	}
	if req.FlatDiscount != nil {
		existingInternship.FlatDiscount = *req.FlatDiscount
	}
	if req.PercentageDiscount != nil {
		existingInternship.PercentageDiscount = *req.PercentageDiscount
	}
	if req.Skills != nil {
		existingInternship.Skills = req.Skills
	}
	if req.LearningOutcomes != nil {
		existingInternship.LearningOutcomes = req.LearningOutcomes
	}
	if req.Prerequisites != nil {
		existingInternship.Prerequisites = req.Prerequisites
	}
	if req.Benefits != nil {
		existingInternship.Benefits = req.Benefits
	}

	existingInternship.UpdatedBy = authContext.UserID

	// Save to repository
	err = s.internshipRepo.Update(ctx, existingInternship)
	if err != nil {
		s.logger.Errorw("Failed to update internship", "error", err, "user_id", authContext.UserID)
		return nil, err
	}

	s.logger.Infow("Internship updated successfully", "internship_id", id, "user_id", authContext.UserID)
	return existingInternship, nil
}

func (s *internshipService) Delete(ctx context.Context, id string, authContext *domainAuth.AuthContext) error {
	// Get existing internship
	existingInternship, err := s.internshipRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if existingInternship == nil {
		return ierr.ErrNotFound
	}

	// Create authorization request
	authRequest := &domainAuth.AccessRequest{
		Subject: authContext,
		Resource: &domainAuth.Resource{
			Type: "internship",
			ID:   id,
			Attributes: map[string]interface{}{
				"created_by": existingInternship.CreatedBy,
				"status":     existingInternship.Status,
			},
		},
		Action: domainAuth.PermissionDeleteInternship,
	}

	// Check authorization
	allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
	if err != nil {
		s.logger.Errorw("Authorization check failed", "error", err, "user_id", authContext.UserID)
		return err
	}

	if !allowed {
		s.logger.Warnw("User not authorized to delete internship", "user_id", authContext.UserID, "internship_id", id)
		return ierr.ErrPermissionDenied
	}

	// Delete from repository (using ID instead of entity)
	err = s.internshipRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to delete internship", "error", err, "user_id", authContext.UserID)
		return err
	}

	s.logger.Infow("Internship deleted successfully", "internship_id", id, "user_id", authContext.UserID)
	return nil
}

func (s *internshipService) List(ctx context.Context, filter *types.InternshipFilter, authContext *domainAuth.AuthContext) ([]*domainInternship.Internship, error) {
	// For listing, we might want to filter based on user role and permissions
	// Students should only see internships they're enrolled in or public ones
	// Instructors can see their own internships
	// Admins can see all internships

	// Create authorization request for listing
	authRequest := &domainAuth.AccessRequest{
		Subject: authContext,
		Resource: &domainAuth.Resource{
			Type: "internship",
		},
		Action: domainAuth.PermissionViewInternship,
	}

	// Check basic permission to view internships
	allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
	if err != nil {
		s.logger.Errorw("Authorization check failed for listing", "error", err, "user_id", authContext.UserID)
		return nil, err
	}

	if !allowed {
		s.logger.Warnw("User not authorized to list internships", "user_id", authContext.UserID)
		return nil, ierr.ErrPermissionDenied
	}

	// Modify filter based on user role
	if authContext.Role == types.UserRoleStudent {
		// Students can only see published internships or ones they're enrolled in
		publishedStatus := types.StatusPublished
		filter.Status = &publishedStatus
		// You might want to add enrollment filtering here
	} else if authContext.Role == types.UserRoleInstructor {
		// Instructors can see their own internships and published ones
		// Note: You'll need to add CreatedBy field to InternshipFilter if it doesn't exist
		// For now, we'll just allow them to see all published internships
		if filter.Status == nil {
			publishedStatus := types.StatusPublished
			filter.Status = &publishedStatus
		}
	}
	// Admins can see all internships (no additional filtering)

	internships, err := s.internshipRepo.List(ctx, filter)
	if err != nil {
		s.logger.Errorw("Failed to list internships", "error", err, "user_id", authContext.UserID)
		return nil, err
	}

	return internships, nil
}
