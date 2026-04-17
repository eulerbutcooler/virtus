package pagination

import (
	"math"
	"net/http"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Page struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Parses `page` and `limit` query params with some defaults.
func FromRequest(r *http.Request) Page {
	page := queryInt(r, "page", DefaultPage)
	limit := queryInt(r, "limit", DefaultLimit)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Page{Page: page, Limit: limit}
}

// Fills in the computed total fields.
func (p Page) WithTotal(total int) Page {
	p.TotalItems = total
	p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.HasNext = p.Page < p.TotalPages
	p.HasPrev = p.Page > 1
	return p
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
