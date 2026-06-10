package handler

import (
	"errors"
	"net/http"

	"github.com/example/ems/internal/auth/dto"
	"github.com/example/ems/internal/auth/service"
	"github.com/example/ems/pkg/apperror"
	"github.com/example/ems/pkg/response"
	"github.com/example/ems/pkg/validation"
	"github.com/gin-gonic/gin"
)

// Handler exposes the auth HTTP endpoints.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a user account and returns an access/refresh token pair.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterRequest  true  "Registration payload"
// @Success      201      {object}  response.Success{data=dto.AuthResponse}
// @Failure      400      {object}  response.Error
// @Failure      409      {object}  response.Error
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}
	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Created(c, "registration successful", res)
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Validates credentials and returns an access/refresh token pair.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginRequest  true  "Login payload"
// @Success      200      {object}  response.Success{data=dto.AuthResponse}
// @Failure      400      {object}  response.Error
// @Failure      401      {object}  response.Error
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "login successful", res)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchanges a valid refresh token for a new token pair.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RefreshRequest  true  "Refresh payload"
// @Success      200      {object}  response.Success{data=dto.TokenResponse}
// @Failure      400      {object}  response.Error
// @Failure      401      {object}  response.Error
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "token refreshed", res)
}

func respondBindError(c *gin.Context, err error) {
	if fields := validation.FieldErrors(err); fields != nil {
		response.FailWithFields(c, http.StatusBadRequest, "validation failed", fields)
		return
	}
	response.Fail(c, http.StatusBadRequest, "invalid request payload")
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperror.ErrConflict):
		response.Fail(c, http.StatusConflict, "email already registered")
	case errors.Is(err, apperror.ErrUnauthorized):
		response.Fail(c, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, apperror.ErrInvalidInput):
		response.Fail(c, http.StatusBadRequest, "invalid input")
	default:
		response.Fail(c, http.StatusInternalServerError, "internal server error")
	}
}
