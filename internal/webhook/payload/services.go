package payload

import "github.com/omkar273/codegeeky/internal/service"

// Services container for all services needed by payload builders
type Services struct {
	UserService       service.UserService
	AuthService       service.AuthService
	CategoryService   service.CategoryService
	OnboardingService service.OnboardingService
	InternshipService service.InternshipService
}

// NewServices creates a new Services container
func NewServices(
	userService service.UserService,
	authService service.AuthService,
	categoryService service.CategoryService,
	onboardingService service.OnboardingService,
	internshipService service.InternshipService,
) *Services {
	return &Services{
		UserService:       userService,
		AuthService:       authService,
		CategoryService:   categoryService,
		OnboardingService: onboardingService,
		InternshipService: internshipService,
	}
}
