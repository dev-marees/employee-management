package service

import (
	"context"

	"github.com/example/ems/internal/dashboard/repository"
)

// StatsResponse is the dashboard payload returned to the frontend.
type StatsResponse struct {
	TotalEmployees    int64                        `json:"total_employees"`
	ActiveEmployees   int64                        `json:"active_employees"`
	InactiveEmployees int64                        `json:"inactive_employees"`
	DepartmentWise    []repository.DepartmentCount `json:"department_wise_count"`
}

// Service defines the dashboard use-cases.
type Service interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStats(ctx context.Context) (*StatsResponse, error) {
	stats, err := s.repo.Stats(ctx)
	if err != nil {
		return nil, err
	}
	deptWise := stats.DepartmentWise
	if deptWise == nil {
		deptWise = []repository.DepartmentCount{}
	}
	return &StatsResponse{
		TotalEmployees:    stats.TotalEmployees,
		ActiveEmployees:   stats.ActiveEmployees,
		InactiveEmployees: stats.InactiveEmployees,
		DepartmentWise:    deptWise,
	}, nil
}
