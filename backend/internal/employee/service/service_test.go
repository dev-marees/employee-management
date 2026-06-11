package service

import (
	"context"
	"testing"

	"github.com/example/ems/internal/employee/dto"
	"github.com/example/ems/internal/employee/model"
	"github.com/example/ems/internal/employee/repository"
	"github.com/example/ems/pkg/apperror"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRepo is an in-memory implementation of repository.Repository.
type fakeRepo struct {
	store map[uuid.UUID]*model.Employee
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{store: map[uuid.UUID]*model.Employee{}}
}

func (f *fakeRepo) Create(_ context.Context, e *model.Employee) error {
	f.store[e.ID] = e
	return nil
}

func (f *fakeRepo) Update(_ context.Context, e *model.Employee) error {
	f.store[e.ID] = e
	return nil
}

func (f *fakeRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	if _, ok := f.store[id]; !ok {
		return apperror.ErrNotFound
	}
	delete(f.store, id)
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id uuid.UUID) (*model.Employee, error) {
	if e, ok := f.store[id]; ok {
		return e, nil
	}
	return nil, apperror.ErrNotFound
}

func (f *fakeRepo) FindByUserID(_ context.Context, userID uuid.UUID) (*model.Employee, error) {
	for _, e := range f.store {
		if e.UserID != nil && *e.UserID == userID {
			return e, nil
		}
	}
	return nil, apperror.ErrNotFound
}

func (f *fakeRepo) List(_ context.Context, _ repository.Filter) ([]model.Employee, int64, error) {
	items := make([]model.Employee, 0, len(f.store))
	for _, e := range f.store {
		items = append(items, *e)
	}
	return items, int64(len(items)), nil
}

func (f *fakeRepo) ExistsByEmail(_ context.Context, email string, excludeID uuid.UUID) (bool, error) {
	for _, e := range f.store {
		if e.Email == email && e.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeRepo) ExistsByCode(_ context.Context, code string) (bool, error) {
	for _, e := range f.store {
		if e.EmployeeCode == code {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeRepo) ExistsByUserID(_ context.Context, userID uuid.UUID) (bool, error) {
	for _, e := range f.store {
		if e.UserID != nil && *e.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func validCreateReq() dto.CreateEmployeeRequest {
	return dto.CreateEmployeeRequest{
		EmployeeCode: "EMP-001",
		FirstName:    "John",
		LastName:     "Smith",
		Email:        "John.Smith@Acme.com",
		Department:   "Engineering",
		Salary:       90000,
		JoiningDate:  "2023-04-01",
	}
}

func TestCreateSuccess(t *testing.T) {
	svc := New(newFakeRepo())
	res, err := svc.Create(context.Background(), validCreateReq())
	require.NoError(t, err)
	assert.Equal(t, "EMP-001", res.EmployeeCode)
	assert.Equal(t, "john.smith@acme.com", res.Email, "email should be normalized")
	assert.Equal(t, "active", res.Status, "status should default to active")
	assert.Equal(t, "2023-04-01", res.JoiningDate)
}

func TestCreateInvalidDate(t *testing.T) {
	svc := New(newFakeRepo())
	req := validCreateReq()
	req.JoiningDate = "01-04-2023" // wrong format
	_, err := svc.Create(context.Background(), req)
	assert.ErrorIs(t, err, apperror.ErrInvalidInput)
}

func TestCreateDuplicateCode(t *testing.T) {
	svc := New(newFakeRepo())
	_, err := svc.Create(context.Background(), validCreateReq())
	require.NoError(t, err)

	dup := validCreateReq()
	dup.Email = "other@acme.com"
	_, err = svc.Create(context.Background(), dup)
	assert.ErrorIs(t, err, apperror.ErrConflict)
}

func TestUpdatePartialFields(t *testing.T) {
	svc := New(newFakeRepo())
	created, err := svc.Create(context.Background(), validCreateReq())
	require.NoError(t, err)

	newSalary := 120000.0
	newStatus := "inactive"
	updated, err := svc.Update(context.Background(), created.ID, dto.UpdateEmployeeRequest{
		Salary: &newSalary,
		Status: &newStatus,
	})
	require.NoError(t, err)
	assert.Equal(t, 120000.0, updated.Salary)
	assert.Equal(t, "inactive", updated.Status)
	assert.Equal(t, "John", updated.FirstName, "unset fields must remain unchanged")
}

func TestUpdateNotFound(t *testing.T) {
	svc := New(newFakeRepo())
	_, err := svc.Update(context.Background(), uuid.New(), dto.UpdateEmployeeRequest{})
	assert.ErrorIs(t, err, apperror.ErrNotFound)
}

func TestDeleteSuccessAndNotFound(t *testing.T) {
	svc := New(newFakeRepo())
	created, err := svc.Create(context.Background(), validCreateReq())
	require.NoError(t, err)

	require.NoError(t, svc.Delete(context.Background(), created.ID))
	assert.ErrorIs(t, svc.Delete(context.Background(), created.ID), apperror.ErrNotFound)
}

func TestListNormalizesPagination(t *testing.T) {
	svc := New(newFakeRepo())
	_, err := svc.Create(context.Background(), validCreateReq())
	require.NoError(t, err)

	items, total, err := svc.List(context.Background(), dto.ListQuery{})
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(1), total)
}
