package httputil

import (
	"fmt"
	"net/http"
)

type Pagination struct {
	Page    int             `json:"page,omitempty"`
	PerPage int             `json:"per_page,omitempty"`
	Link    PaginationLinks `json:"links,omitempty"`
}

type PaginationLinks struct {
	Prev string `json:"prev,omitempty"`
	Next string `json:"next,omitempty"`
}

func NewPagination(size int, from int) Pagination {
	return Pagination{
		Page:    computePage(from, size),
		PerPage: size,
	}
}

func (p *Pagination) SetLinks(r *http.Request) {
	links := PaginationLinks{
		Next: buildURLWithPagination(r, *p, p.PerPage),
	}

	if p.Page > 1 {
		links.Prev = buildURLWithPagination(r, *p, -p.PerPage)
	}

	p.Link = links
}

func computeOffset(page int, size int) int {
	return (page - 1) * size
}

func computePage(from int, size int) int {
	return (from / size) + 1
}

func buildURLWithQuery(r *http.Request, values map[string]int) string {
	copy := r.URL.Query()
	for p, v := range values {
		copy.Set(p, fmt.Sprint(v))
	}
	r.URL.RawQuery = copy.Encode()
	return r.URL.String()
}

func buildURLWithPagination(r *http.Request, p Pagination, delta int) string {
	v := map[string]int{
		"size": p.PerPage,
		"from": computeOffset(p.Page, p.PerPage) + delta,
	}
	return buildURLWithQuery(r, v)
}
