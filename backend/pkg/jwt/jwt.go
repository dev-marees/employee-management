// Package jwt issues and verifies access/refresh tokens used for authentication.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType distinguishes access tokens from refresh tokens so a refresh token
// cannot be replayed as an access token (and vice-versa).
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrWrongType    = errors.New("unexpected token type")
)

// Claims is the JWT payload carried in both token types.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// Manager issues and validates tokens. Access and refresh tokens are signed
// with independent secrets.
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
}

func NewManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, issuer string) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		issuer:        issuer,
	}
}

// TokenPair is returned on login/refresh.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // access token lifetime in seconds
}

// GeneratePair issues a fresh access+refresh token pair for a user.
func (m *Manager) GeneratePair(userID uuid.UUID, email, role string) (*TokenPair, error) {
	access, err := m.sign(userID, email, role, AccessToken, m.accessSecret, m.accessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := m.sign(userID, email, role, RefreshToken, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

func (m *Manager) sign(userID uuid.UUID, email, role string, typ TokenType, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        uuid.NewString(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// VerifyAccess validates an access token and returns its claims.
func (m *Manager) VerifyAccess(token string) (*Claims, error) {
	return m.verify(token, m.accessSecret, AccessToken)
}

// VerifyRefresh validates a refresh token and returns its claims.
func (m *Manager) VerifyRefresh(token string) (*Claims, error) {
	return m.verify(token, m.refreshSecret, RefreshToken)
}

func (m *Manager) verify(token string, secret []byte, expected TokenType) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Type != expected {
		return nil, ErrWrongType
	}
	return claims, nil
}
