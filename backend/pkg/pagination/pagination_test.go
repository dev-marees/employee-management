package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAppliesDefaults(t *testing.T) {
	p := Params{}
	p.Normalize()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.Limit)
}

func TestNormalizeClampsLimit(t *testing.T) {
	p := Params{Page: -3, Limit: 5000}
	p.Normalize()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, maxLimit, p.Limit)
}

func TestOffset(t *testing.T) {
	p := Params{Page: 3, Limit: 10}
	assert.Equal(t, 20, p.Offset())
}

func TestNewResultComputesTotalPages(t *testing.T) {
	tests := []struct {
		total int64
		limit int
		want  int
	}{
		{total: 100, limit: 10, want: 10},
		{total: 101, limit: 10, want: 11},
		{total: 0, limit: 10, want: 0},
		{total: 5, limit: 10, want: 1},
	}
	for _, tc := range tests {
		res := NewResult(nil, Params{Page: 1, Limit: tc.limit}, tc.total)
		assert.Equal(t, tc.want, res.TotalPages, "total=%d limit=%d", tc.total, tc.limit)
	}
}
