package dto

// OnboardingRequest is the request for the onboarding service
// It extends the SignupRequest to include the role
type OnboardingRequest struct {
	SignupRequest
}

func (r *OnboardingRequest) Validate() error {

	// validate signup request
	if err := r.SignupRequest.Validate(); err != nil {
		return err
	}

	return nil
}

// OnboardingResponse is the response for the onboarding service
// It extends the SignupResponse to include the role
type OnboardingResponse struct {
	SignupResponse
}
