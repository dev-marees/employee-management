package dto

import (
	"time"

	"github.com/example/ems/internal/auth/model"
	"github.com/google/uuid"
)

// RegisterRequest is the payload for POST /auth/register.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120" example:"Jane Doe"`
	Email    string `json:"email" binding:"required,email" example:"jane@acme.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"S3cretPass"`
	Role     string `json:"role" binding:"omitempty,oneof=Admin HR Employee" example:"HR"`
}

// LoginRequest is the payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"jane@acme.com"`
	Password string `json:"password" binding:"required" example:"S3cretPass"`
}

// RefreshRequest is the payload for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest is the payload for POST /auth/change-password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"Temp@1234"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=72" example:"MyNewPass1"`
}

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
}

// AuthResponse bundles the issued tokens with the user profile.
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
}

// TokenResponse is returned by the refresh endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ToUserResponse maps a User model to its public DTO.
func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:                 u.ID,
		Name:               u.Name,
		Email:              u.Email,
		Role:               string(u.Role),
		MustChangePassword: u.MustChangePassword,
		CreatedAt:          u.CreatedAt,
	}
}
