package handler

import (
	"net/http"

	"github.com/example/ems/internal/dashboard/service"
	"github.com/example/ems/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler exposes the dashboard HTTP endpoint.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// Stats godoc
// @Summary      Dashboard statistics
// @Description  Returns aggregate employee counts and a department breakdown.
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Success{data=service.StatsResponse}
// @Failure      401  {object}  response.Error
// @Router       /dashboard [get]
func (h *Handler) Stats(c *gin.Context) {
	res, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to load dashboard")
		return
	}
	response.OK(c, "dashboard fetched", res)
}
