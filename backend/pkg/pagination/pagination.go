// Package pagination provides reusable request/response helpers for list endpoints.
package pagination

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

// Params holds normalized pagination input.
type Params struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

// Normalize clamps page/limit into safe bounds and applies defaults.
func (p *Params) Normalize() {
	if p.Page < 1 {
		p.Page = defaultPage
	}
	if p.Limit < 1 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
}

// Offset is the SQL offset for the current page.
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Result is the standard paginated list envelope returned to the frontend.
type Result struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// NewResult builds a Result, computing total_pages from total and limit.
func NewResult(data interface{}, params Params, total int64) Result {
	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}
	return Result{
		Data:       data,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
