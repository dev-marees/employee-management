// Package response provides a consistent JSON envelope for all API responses,
// matching the success/error contract expected by the frontend.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success is the standard success envelope.
type Success struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// Error is the standard error envelope. Errors carries field-level validation
// details (field -> message) when applicable.
type Error struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// OK writes a 200 response with data.
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Success{Success: true, Message: message, Data: data})
}

// Created writes a 201 response with data.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Success{Success: true, Message: message, Data: data})
}

// JSON writes an arbitrary success status with data.
func JSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Success{Success: true, Message: message, Data: data})
}

// Fail writes an error response with the given status and message.
func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Error{Success: false, Message: message})
}

// FailWithFields writes an error response carrying field-level errors.
func FailWithFields(c *gin.Context, status int, message string, fields map[string]string) {
	c.JSON(status, Error{Success: false, Message: message, Errors: fields})
}
