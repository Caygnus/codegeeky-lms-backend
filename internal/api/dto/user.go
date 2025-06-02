package dto

type MeResponse struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Role     string `json:"role,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
