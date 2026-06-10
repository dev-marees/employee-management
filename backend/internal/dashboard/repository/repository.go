package repository

import (
	"context"

	"github.com/example/ems/internal/employee/model"
	"gorm.io/gorm"
)

// DepartmentCount is one row of the department-wise aggregation.
type DepartmentCount struct {
	Department string `json:"department"`
	Count      int64  `json:"count"`
}

// Stats is the raw aggregation produced by the repository.
type Stats struct {
	TotalEmployees    int64
	ActiveEmployees   int64
	InactiveEmployees int64
	DepartmentWise    []DepartmentCount
}

// Repository abstracts dashboard aggregation queries.
type Repository interface {
	Stats(ctx context.Context) (*Stats, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Stats(ctx context.Context) (*Stats, error) {
	db := r.db.WithContext(ctx).Model(&model.Employee{})

	var total, active, inactive int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Session(&gorm.Session{}).Where("status = ?", model.StatusActive).Count(&active).Error; err != nil {
		return nil, err
	}
	if err := db.Session(&gorm.Session{}).Where("status = ?", model.StatusInactive).Count(&inactive).Error; err != nil {
		return nil, err
	}

	var deptWise []DepartmentCount
	err := db.Session(&gorm.Session{}).
		Select("department, COUNT(*) AS count").
		Group("department").
		Order("count DESC").
		Scan(&deptWise).Error
	if err != nil {
		return nil, err
	}

	return &Stats{
		TotalEmployees:    total,
		ActiveEmployees:   active,
		InactiveEmployees: inactive,
		DepartmentWise:    deptWise,
	}, nil
}
