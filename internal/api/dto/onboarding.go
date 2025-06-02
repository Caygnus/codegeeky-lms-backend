package dto

import ierr "github.com/omkar273/codegeeky/internal/errors"

// OnboardingRequest is the request for the onboarding service
// It extends the SignupRequest to include the role
type OnboardingRequest struct {
	SignupRequest
	ProviderUserID string `json:"provider_user_id" validate:"required"`
}

func (r *OnboardingRequest) Validate() error {

	// validate signup request
	if err := r.SignupRequest.Validate(); err != nil {
		return err
	}

	// validate provider user id
	if r.ProviderUserID == "" {
		return ierr.NewError("provider user id is required").
			Mark(ierr.ErrValidation)
	}

	return nil
}

// OnboardingResponse is the response for the onboarding service
// It extends the SignupResponse to include the role
type OnboardingResponse struct {
	SignupResponse
}
