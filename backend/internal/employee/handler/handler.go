package handler

import (
	"errors"
	"net/http"

	"github.com/example/ems/internal/employee/dto"
	"github.com/example/ems/internal/employee/service"
	"github.com/example/ems/internal/middleware"
	"github.com/example/ems/pkg/apperror"
	"github.com/example/ems/pkg/pagination"
	"github.com/example/ems/pkg/response"
	"github.com/example/ems/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler exposes the employee HTTP endpoints.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc
// @Summary      List employees
// @Description  Returns a paginated, searchable, filterable and sortable list.
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        page            query     int     false  "Page number"        default(1)
// @Param        limit           query     int     false  "Items per page"     default(10)
// @Param        search          query     string  false  "Search code/name/email/department"
// @Param        department      query     string  false  "Filter by department"
// @Param        status          query     string  false  "Filter by status (active|inactive)"
// @Param        joined_from     query     string  false  "Joined date from (YYYY-MM-DD)"
// @Param        joined_to       query     string  false  "Joined date to (YYYY-MM-DD)"
// @Param        sort_by         query     string  false  "Sort by (name|salary|joining_date)"
// @Param        sort_direction  query     string  false  "Sort direction (asc|desc)"
// @Success      200             {object}  response.Success{data=pagination.Result}
// @Failure      401             {object}  response.Error
// @Router       /employees [get]
func (h *Handler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		respondBindError(c, err)
		return
	}
	items, total, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	result := pagination.NewResult(items, q.Params, total)
	response.OK(c, "employees fetched", result)
}

// Me godoc
// @Summary      Get the current user's employee record
// @Description  Returns the employee profile linked to the authenticated user.
//
//	Any authenticated role may call this; it only ever returns the caller's
//	own record. Returns 404 if no employee is linked to the user.
//
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Success{data=dto.EmployeeResponse}
// @Failure      401  {object}  response.Error
// @Failure      404  {object}  response.Error
// @Router       /me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthenticated")
		return
	}
	res, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "employee fetched", res)
}

// Get godoc
// @Summary      Get employee by ID
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Employee ID (uuid)"
// @Success      200  {object}  response.Success{data=dto.EmployeeResponse}
// @Failure      404  {object}  response.Error
// @Router       /employees/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "employee fetched", res)
}

// Create godoc
// @Summary      Create employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      dto.CreateEmployeeRequest  true  "Employee payload"
// @Success      201      {object}  response.Success{data=dto.EmployeeResponse}
// @Failure      400      {object}  response.Error
// @Failure      409      {object}  response.Error
// @Router       /employees [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}
	res, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Created(c, "employee created successfully", res)
}

// Update godoc
// @Summary      Update employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                     true  "Employee ID (uuid)"
// @Param        payload  body      dto.UpdateEmployeeRequest  true  "Fields to update"
// @Success      200      {object}  response.Success{data=dto.EmployeeResponse}
// @Failure      400      {object}  response.Error
// @Failure      404      {object}  response.Error
// @Router       /employees/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}
	res, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "employee updated successfully", res)
}

// Delete godoc
// @Summary      Delete employee (soft delete)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Employee ID (uuid)"
// @Success      200  {object}  response.Success
// @Failure      404  {object}  response.Error
// @Router       /employees/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, "employee deleted successfully", nil)
}

func parseID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid employee id")
		return uuid.Nil, false
	}
	return id, true
}

func respondBindError(c *gin.Context, err error) {
	if fields := validation.FieldErrors(err); fields != nil {
		response.FailWithFields(c, http.StatusBadRequest, "validation failed", fields)
		return
	}
	response.Fail(c, http.StatusBadRequest, "invalid request payload")
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		response.Fail(c, http.StatusNotFound, "employee not found")
	case errors.Is(err, apperror.ErrConflict):
		response.Fail(c, http.StatusConflict, "employee code or email already exists")
	case errors.Is(err, apperror.ErrInvalidInput):
		response.Fail(c, http.StatusBadRequest, "invalid input (check joining_date format YYYY-MM-DD)")
	default:
		response.Fail(c, http.StatusInternalServerError, "internal server error")
	}
}
