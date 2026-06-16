package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	authmodel "github.com/example/ems/internal/auth/model"
	"github.com/example/ems/internal/employee/model"
	"github.com/example/ems/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// empCodeLockKey is an arbitrary, stable key for the Postgres transaction-level
// advisory lock that serializes employee-code generation across concurrent
// creates (prevents two requests from picking the same EMP### number).
const empCodeLockKey int64 = 480127

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
	// CreateWithUser provisions a user account and its employee record in one
	// transaction, auto-generating the employee code and linking user_id.
	CreateWithUser(ctx context.Context, e *model.Employee, u *authmodel.User) error
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

func (r *repository) CreateWithUser(ctx context.Context, e *model.Employee, u *authmodel.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Serialize code generation so concurrent creates can't pick the same
		// number. The lock is held until the transaction commits/rolls back.
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", empCodeLockKey).Error; err != nil {
			return err
		}

		// Email must be free across both user accounts and employee records.
		var cnt int64
		if err := tx.Model(&authmodel.User{}).Where("email = ?", u.Email).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			return apperror.ErrConflict
		}
		if err := tx.Model(&model.Employee{}).Where("email = ?", e.Email).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			return apperror.ErrConflict
		}

		if err := tx.Create(u).Error; err != nil {
			return err
		}

		code, err := nextEmployeeCode(tx)
		if err != nil {
			return err
		}
		e.EmployeeCode = code
		e.UserID = &u.ID
		return tx.Create(e).Error
	})
}

// nextEmployeeCode returns the next sequential code in EMP### format (EMP001,
// EMP002, …). It reads the highest existing numeric suffix — including
// soft-deleted rows (Unscoped) so codes are never reused — and increments it.
// Must be called inside the advisory-locked transaction.
func nextEmployeeCode(tx *gorm.DB) (string, error) {
	var maxN int64
	row := tx.Unscoped().Model(&model.Employee{}).
		Where("employee_code ~ ?", `^EMP[0-9]+$`).
		Select(`COALESCE(MAX(CAST(SUBSTRING(employee_code FROM 4) AS INTEGER)), 0)`).
		Row()
	if err := row.Scan(&maxN); err != nil {
		return "", err
	}
	return fmt.Sprintf("EMP%03d", maxN+1), nil
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
