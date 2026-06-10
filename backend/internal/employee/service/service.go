package service

import (
	"context"
	"strings"
	"time"

	"github.com/example/ems/internal/employee/dto"
	"github.com/example/ems/internal/employee/model"
	"github.com/example/ems/internal/employee/repository"
	"github.com/example/ems/pkg/apperror"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

// Service defines the employee use-cases.
type Service interface {
	Create(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.EmployeeResponse, error)
	List(ctx context.Context, q dto.ListQuery) ([]dto.EmployeeResponse, int64, error)
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	joiningDate, err := time.Parse(dateLayout, req.JoiningDate)
	if err != nil {
		return nil, apperror.ErrInvalidInput
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if exists, err := s.repo.ExistsByCode(ctx, req.EmployeeCode); err != nil {
		return nil, err
	} else if exists {
		return nil, apperror.ErrConflict
	}
	if exists, err := s.repo.ExistsByEmail(ctx, email, uuid.Nil); err != nil {
		return nil, err
	} else if exists {
		return nil, apperror.ErrConflict
	}

	status := model.Status(req.Status)
	if req.Status == "" {
		status = model.StatusActive
	}

	e := &model.Employee{
		ID:           uuid.New(),
		EmployeeCode: req.EmployeeCode,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        email,
		Phone:        req.Phone,
		Department:   req.Department,
		Designation:  req.Designation,
		Salary:       req.Salary,
		JoiningDate:  joiningDate,
		Status:       status,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	resp := dto.ToResponse(e)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		e.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		e.LastName = *req.LastName
	}
	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		if exists, err := s.repo.ExistsByEmail(ctx, email, e.ID); err != nil {
			return nil, err
		} else if exists {
			return nil, apperror.ErrConflict
		}
		e.Email = email
	}
	if req.Phone != nil {
		e.Phone = *req.Phone
	}
	if req.Department != nil {
		e.Department = *req.Department
	}
	if req.Designation != nil {
		e.Designation = *req.Designation
	}
	if req.Salary != nil {
		e.Salary = *req.Salary
	}
	if req.JoiningDate != nil {
		jd, err := time.Parse(dateLayout, *req.JoiningDate)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		e.JoiningDate = jd
	}
	if req.Status != nil {
		e.Status = model.Status(*req.Status)
	}

	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	resp := dto.ToResponse(e)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*dto.EmployeeResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToResponse(e)
	return &resp, nil
}

func (s *service) List(ctx context.Context, q dto.ListQuery) ([]dto.EmployeeResponse, int64, error) {
	q.Params.Normalize()

	filter := repository.Filter{
		Search:        strings.TrimSpace(q.Search),
		Department:    strings.TrimSpace(q.Department),
		Status:        strings.TrimSpace(q.Status),
		JoinedFrom:    strings.TrimSpace(q.JoinedFrom),
		JoinedTo:      strings.TrimSpace(q.JoinedTo),
		SortBy:        q.SortBy,
		SortDirection: q.SortDirection,
		Limit:         q.Limit,
		Offset:        q.Offset(),
	}

	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToResponseList(items), total, nil
}
