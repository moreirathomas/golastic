package httputil

import (
	"fmt"
	"net/http"
	"net/url"
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

func NewPagination(size, from int) Pagination {
	return Pagination{
		Page:    computePage(from, size),
		PerPage: size,
	}
}

func (p *Pagination) SetLinks(r *http.Request, total int) {
	links := PaginationLinks{}

	if p.Page > 1 {
		links.Prev = buildURLWithPagination(r, *p, -p.PerPage)
	}
	if total > p.Page*p.PerPage {
		links.Next = buildURLWithPagination(r, *p, p.PerPage)
	}

	p.Link = links
}

func computeOffset(page, size int) int {
	return (page - 1) * size
}

func computePage(from, size int) int {
	return (from / size) + 1
}

func getBaseURL(r *http.Request) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.URL.Path,
	}
}

func buildURLWithQuery(r *http.Request, values map[string]int) string {
	newURL := getBaseURL(r)
	query := r.URL.Query()
	for p, v := range values {
		query.Set(p, fmt.Sprint(v))
	}
	newURL.RawQuery = query.Encode()
	return newURL.String()
}

func buildURLWithPagination(r *http.Request, p Pagination, delta int) string {
	v := map[string]int{
		"size": p.PerPage,
		"from": computeOffset(p.Page, p.PerPage) + delta,
	}
	return buildURLWithQuery(r, v)
}
