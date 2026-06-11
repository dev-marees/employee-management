package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status enumerates an employee's employment status.
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func (s Status) Valid() bool {
	return s == StatusActive || s == StatusInactive
}

// Employee is the core domain entity. Soft deletes are enabled via DeletedAt.
type Employee struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	// UserID links this employee record to its authenticatable user account.
	// Nullable: an employee may exist before a login is provisioned. The unique
	// index enforces at most one employee per user.
	UserID       *uuid.UUID     `gorm:"type:uuid;uniqueIndex" json:"user_id,omitempty"`
	EmployeeCode string         `gorm:"type:varchar(40);uniqueIndex;not null" json:"employee_code"`
	FirstName    string         `gorm:"type:varchar(80);not null;index" json:"first_name"`
	LastName     string         `gorm:"type:varchar(80);not null;index" json:"last_name"`
	Email        string         `gorm:"type:varchar(160);uniqueIndex;not null" json:"email"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone"`
	Department   string         `gorm:"type:varchar(80);index" json:"department"`
	Designation  string         `gorm:"type:varchar(80)" json:"designation"`
	Salary       float64        `gorm:"type:numeric(12,2);not null;default:0" json:"salary"`
	JoiningDate  time.Time      `gorm:"type:date;not null;index" json:"joining_date"`
	Status       Status         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Employee) TableName() string { return "employees" }

func (e *Employee) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
