package dto

import (
	"time"

	"github.com/example/ems/internal/employee/model"
	"github.com/example/ems/pkg/pagination"
	"github.com/google/uuid"
)

// CreateEmployeeRequest is the payload for POST /employees. The backend
// provisions a linked user account (using Email + Password) and auto-generates
// the employee code, so neither employee_code nor user_id is accepted here.
type CreateEmployeeRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=1,max=80" example:"John"`
	LastName    string  `json:"last_name" binding:"required,min=1,max=80" example:"Smith"`
	Email       string  `json:"email" binding:"required,email" example:"john.smith@acme.com"`
	// Password is the temporary password for the new user account. The user is
	// forced to change it on first login.
	Password    string  `json:"password" binding:"required,min=8,max=72" example:"Temp@1234"`
	Phone       string  `json:"phone" binding:"omitempty,max=20" example:"+14155550123"`
	Department  string  `json:"department" binding:"required,max=80" example:"Engineering"`
	Designation string  `json:"designation" binding:"omitempty,max=80" example:"Senior Engineer"`
	Salary      float64 `json:"salary" binding:"gte=0" example:"95000"`
	JoiningDate string  `json:"joining_date" binding:"required" example:"2023-04-01"` // YYYY-MM-DD
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive" example:"active"`
}

// UpdateEmployeeRequest is the payload for PUT /employees/:id. All fields are
// optional; only provided (non-nil) fields are applied.
type UpdateEmployeeRequest struct {
	FirstName   *string  `json:"first_name" binding:"omitempty,min=1,max=80"`
	LastName    *string  `json:"last_name" binding:"omitempty,min=1,max=80"`
	Email       *string  `json:"email" binding:"omitempty,email"`
	Phone       *string  `json:"phone" binding:"omitempty,max=20"`
	Department  *string  `json:"department" binding:"omitempty,max=80"`
	Designation *string  `json:"designation" binding:"omitempty,max=80"`
	Salary      *float64 `json:"salary" binding:"omitempty,gte=0"`
	JoiningDate *string  `json:"joining_date" binding:"omitempty"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListQuery captures search/filter/sort/pagination query parameters for
// GET /employees.
type ListQuery struct {
	pagination.Params
	Search        string `form:"search"`         // matches code, name, email, department
	Department    string `form:"department"`     // exact department filter
	Status        string `form:"status"`         // active|inactive
	JoinedFrom    string `form:"joined_from"`    // YYYY-MM-DD inclusive
	JoinedTo      string `form:"joined_to"`      // YYYY-MM-DD inclusive
	SortBy        string `form:"sort_by"`        // name|salary|joining_date
	SortDirection string `form:"sort_direction"` // asc|desc
}

// EmployeeResponse is the public representation of an employee.
type EmployeeResponse struct {
	ID           uuid.UUID  `json:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	EmployeeCode string     `json:"employee_code"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Department   string    `json:"department"`
	Designation  string    `json:"designation"`
	Salary       float64   `json:"salary"`
	JoiningDate  string    `json:"joining_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const dateLayout = "2006-01-02"

// ToResponse maps an Employee model to its public DTO.
func ToResponse(e *model.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:           e.ID,
		UserID:       e.UserID,
		EmployeeCode: e.EmployeeCode,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Email:        e.Email,
		Phone:        e.Phone,
		Department:   e.Department,
		Designation:  e.Designation,
		Salary:       e.Salary,
		JoiningDate:  e.JoiningDate.Format(dateLayout),
		Status:       string(e.Status),
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// ToResponseList maps a slice of Employee models to DTOs.
func ToResponseList(items []model.Employee) []EmployeeResponse {
	out := make([]EmployeeResponse, 0, len(items))
	for i := range items {
		out = append(out, ToResponse(&items[i]))
	}
	return out
}
