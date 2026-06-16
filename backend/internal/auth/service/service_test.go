package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/ems/internal/auth/dto"
	"github.com/example/ems/internal/auth/model"
	"github.com/example/ems/pkg/apperror"
	"github.com/example/ems/pkg/hash"
	"github.com/example/ems/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRepo is an in-memory implementation of repository.Repository for tests.
type fakeRepo struct {
	byEmail map[string]*model.User
	byID    map[uuid.UUID]*model.User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		byEmail: map[string]*model.User{},
		byID:    map[uuid.UUID]*model.User{},
	}
}

func (f *fakeRepo) Create(_ context.Context, u *model.User) error {
	f.byEmail[u.Email] = u
	f.byID[u.ID] = u
	return nil
}

func (f *fakeRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	if u, ok := f.byEmail[email]; ok {
		return u, nil
	}
	return nil, apperror.ErrNotFound
}

func (f *fakeRepo) FindByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	if u, ok := f.byID[id]; ok {
		return u, nil
	}
	return nil, apperror.ErrNotFound
}

func (f *fakeRepo) ExistsByEmail(_ context.Context, email string) (bool, error) {
	_, ok := f.byEmail[email]
	return ok, nil
}

func (f *fakeRepo) UpdatePassword(_ context.Context, id uuid.UUID, passwordHash string) error {
	u, ok := f.byID[id]
	if !ok {
		return apperror.ErrNotFound
	}
	u.PasswordHash = passwordHash
	u.MustChangePassword = false
	return nil
}

func newTestService() (Service, *fakeRepo) {
	repo := newFakeRepo()
	mgr := jwt.NewManager("a", "r", time.Minute, time.Hour, "test")
	return New(repo, mgr), repo
}

func TestRegisterSuccess(t *testing.T) {
	svc, repo := newTestService()

	res, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name:     "Jane Doe",
		Email:    "Jane@Example.com",
		Password: "supersecret",
		Role:     "HR",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.Equal(t, "jane@example.com", res.User.Email, "email should be normalized to lowercase")

	// Password must be stored hashed, never in plaintext.
	stored := repo.byEmail["jane@example.com"]
	require.NotNil(t, stored)
	assert.NotEqual(t, "supersecret", stored.PasswordHash)
	assert.True(t, hash.Compare(stored.PasswordHash, "supersecret"))
}

func TestRegisterDuplicateEmail(t *testing.T) {
	svc, _ := newTestService()
	req := dto.RegisterRequest{Name: "A", Email: "dup@example.com", Password: "supersecret"}

	_, err := svc.Register(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.Register(context.Background(), req)
	assert.ErrorIs(t, err, apperror.ErrConflict)
}

func TestRegisterDefaultsToEmployeeRole(t *testing.T) {
	svc, repo := newTestService()
	_, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name: "No Role", Email: "norole@example.com", Password: "supersecret",
	})
	require.NoError(t, err)
	assert.Equal(t, model.RoleEmployee, repo.byEmail["norole@example.com"].Role)
}

func TestLoginSuccess(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name: "Log In", Email: "login@example.com", Password: "supersecret",
	})
	require.NoError(t, err)

	res, err := svc.Login(context.Background(), dto.LoginRequest{
		Email: "login@example.com", Password: "supersecret",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
}

func TestLoginWrongPassword(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name: "X", Email: "wrong@example.com", Password: "supersecret",
	})
	require.NoError(t, err)

	_, err = svc.Login(context.Background(), dto.LoginRequest{
		Email: "wrong@example.com", Password: "incorrect",
	})
	assert.ErrorIs(t, err, apperror.ErrUnauthorized)
}

func TestLoginUnknownUser(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Email: "ghost@example.com", Password: "whatever",
	})
	assert.ErrorIs(t, err, apperror.ErrUnauthorized)
}

func TestRefreshIssuesNewPair(t *testing.T) {
	svc, _ := newTestService()
	reg, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name: "R", Email: "refresh@example.com", Password: "supersecret",
	})
	require.NoError(t, err)

	res, err := svc.Refresh(context.Background(), reg.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)
}

func TestRefreshRejectsAccessToken(t *testing.T) {
	svc, _ := newTestService()
	reg, err := svc.Register(context.Background(), dto.RegisterRequest{
		Name: "R", Email: "badrefresh@example.com", Password: "supersecret",
	})
	require.NoError(t, err)

	// Passing an access token where a refresh token is expected must fail.
	_, err = svc.Refresh(context.Background(), reg.AccessToken)
	assert.ErrorIs(t, err, apperror.ErrUnauthorized)
}
