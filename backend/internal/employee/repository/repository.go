package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ems/internal/employee/model"
	"github.com/example/ems/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Filter captures the validated, repository-level query criteria.
type Filter struct {
	Search        string
	Department    string
	Status        string
	JoinedFrom    string // YYYY-MM-DD
	JoinedTo      string // YYYY-MM-DD
	SortBy        string // name|salary|joining_date
	SortDirection string // asc|desc
	Limit         int
	Offset        int
}

// allowedSortColumns whitelists sortable columns to prevent SQL injection via
// the sort_by query parameter.
var allowedSortColumns = map[string]string{
	"name":         "first_name",
	"salary":       "salary",
	"joining_date": "joining_date",
}

// Repository abstracts employee persistence.
type Repository interface {
	Create(ctx context.Context, e *model.Employee) error
	Update(ctx context.Context, e *model.Employee) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Employee, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Employee, error)
	List(ctx context.Context, f Filter) ([]model.Employee, int64, error)
	ExistsByEmail(ctx context.Context, email string, excludeID uuid.UUID) (bool, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, e *model.Employee) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *repository) Update(ctx context.Context, e *model.Employee) error {
	return r.db.WithContext(ctx).Save(e).Error
}

func (r *repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Employee{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*model.Employee, error) {
	var e model.Employee
	err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Employee, error) {
	var e model.Employee
	err := r.db.WithContext(ctx).First(&e, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) List(ctx context.Context, f Filter) ([]model.Employee, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Employee{})
	q = applyFilters(q, f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting (whitelisted column + direction).
	column := "created_at"
	if c, ok := allowedSortColumns[f.SortBy]; ok {
		column = c
	}
	direction := "asc"
	if strings.EqualFold(f.SortDirection, "desc") {
		direction = "desc"
	}
	q = q.Order(column + " " + direction)

	if f.Limit > 0 {
		q = q.Limit(f.Limit).Offset(f.Offset)
	}

	var items []model.Employee
	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *repository) ExistsByEmail(ctx context.Context, email string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Employee{}).Where("email = ?", email)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *repository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Employee{}).Where("employee_code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *repository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Employee{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

// applyFilters builds the WHERE clause shared by Count and Find so totals and
// pages stay consistent.
func applyFilters(q *gorm.DB, f Filter) *gorm.DB {
	if f.Search != "" {
		like := "%" + strings.ToLower(f.Search) + "%"
		q = q.Where(
			"LOWER(employee_code) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(department) LIKE ?",
			like, like, like, like, like,
		)
	}
	if f.Department != "" {
		q = q.Where("department = ?", f.Department)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.JoinedFrom != "" {
		q = q.Where("joining_date >= ?", f.JoinedFrom)
	}
	if f.JoinedTo != "" {
		q = q.Where("joining_date <= ?", f.JoinedTo)
	}
	return q
}
