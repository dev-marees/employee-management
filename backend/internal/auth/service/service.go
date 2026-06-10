package service

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ems/internal/auth/dto"
	"github.com/example/ems/internal/auth/model"
	"github.com/example/ems/internal/auth/repository"
	"github.com/example/ems/pkg/apperror"
	"github.com/example/ems/pkg/hash"
	"github.com/example/ems/pkg/jwt"
	"github.com/google/uuid"
)

// Service defines the auth use-cases.
type Service interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
}

type service struct {
	repo repository.Repository
	jwt  *jwt.Manager
}

func New(repo repository.Repository, jwtMgr *jwt.Manager) Service {
	return &service{repo: repo, jwt: jwtMgr}
}

func (s *service) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.ErrConflict
	}

	role := model.Role(req.Role)
	if req.Role == "" {
		role = model.RoleEmployee
	}
	if !role.Valid() {
		return nil, apperror.ErrInvalidInput
	}

	hashed, err := hash.Password(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        email,
		PasswordHash: hashed,
		Role:         role,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *service) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			// Avoid leaking which part failed.
			return nil, apperror.ErrUnauthorized
		}
		return nil, err
	}
	if !hash.Compare(user.PasswordHash, req.Password) {
		return nil, apperror.ErrUnauthorized
	}
	return s.buildAuthResponse(user)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.jwt.VerifyRefresh(refreshToken)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}
	// Confirm the user still exists (and pick up role changes).
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}
	pair, err := s.jwt.GeneratePair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    pair.TokenType,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (s *service) buildAuthResponse(user *model.User) (*dto.AuthResponse, error) {
	pair, err := s.jwt.GeneratePair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    pair.TokenType,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
