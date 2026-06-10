package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() *Manager {
	return NewManager("access-secret", "refresh-secret", time.Minute, time.Hour, "test")
}

func TestGenerateAndVerifyAccess(t *testing.T) {
	m := newTestManager()
	id := uuid.New()

	pair, err := m.GeneratePair(id, "user@example.com", "Admin")
	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)

	claims, err := m.VerifyAccess(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, id, claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, "Admin", claims.Role)
	assert.Equal(t, AccessToken, claims.Type)
}

func TestAccessTokenRejectedAsRefresh(t *testing.T) {
	m := newTestManager()
	pair, err := m.GeneratePair(uuid.New(), "u@e.com", "HR")
	require.NoError(t, err)

	// An access token must not validate as a refresh token. With distinct
	// secrets the signature check fails first; the type check is defense in
	// depth for the shared-secret case. Either way the token is rejected.
	_, err = m.VerifyRefresh(pair.AccessToken)
	assert.Error(t, err)
}

// TestWrongTypeWithSharedSecret exercises the type-mismatch branch directly by
// using the same secret for both token kinds.
func TestWrongTypeWithSharedSecret(t *testing.T) {
	m := NewManager("same", "same", time.Minute, time.Hour, "test")
	pair, err := m.GeneratePair(uuid.New(), "u@e.com", "HR")
	require.NoError(t, err)

	_, err = m.VerifyRefresh(pair.AccessToken)
	assert.ErrorIs(t, err, ErrWrongType)
}

func TestVerifyRejectsTamperedToken(t *testing.T) {
	m := newTestManager()
	pair, err := m.GeneratePair(uuid.New(), "u@e.com", "Employee")
	require.NoError(t, err)

	_, err = m.VerifyAccess(pair.AccessToken + "tampered")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestExpiredTokenRejected(t *testing.T) {
	m := NewManager("a", "r", -time.Minute, time.Hour, "test") // already expired
	pair, err := m.GeneratePair(uuid.New(), "u@e.com", "Employee")
	require.NoError(t, err)

	_, err = m.VerifyAccess(pair.AccessToken)
	assert.ErrorIs(t, err, ErrInvalidToken)
}
