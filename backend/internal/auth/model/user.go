package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role enumerates the application's authorization roles.
type Role string

const (
	RoleAdmin    Role = "Admin"
	RoleHR       Role = "HR"
	RoleEmployee Role = "Employee"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleHR, RoleEmployee:
		return true
	default:
		return false
	}
}

// User is an authenticatable account.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(120);not null" json:"name"`
	Email        string    `gorm:"type:varchar(160);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         Role      `gorm:"type:varchar(20);not null;default:'Employee'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// BeforeCreate guarantees an ID even when the DB default is unavailable.
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
