package pagination

import (
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Links   Links `json:"links"`
}

type Links struct {
	Prev string `json:"prev,omitempty"`
	Next string `json:"next,omitempty"`
}

// PageToOffset returns the pagination offset value
// corresponding to the given page and size.
func PageToOffset(page, size int) int {
	return (page - 1) * size
}

// OffsetToPage returns the pagination page value
// corresponding to the given offset and size.
func OffsetToPage(offset, size int) int {
	return (offset / size) + 1
}

func New(r *http.Request, total, page, size int) (Pagination, error) {
	p := Pagination{
		Page:    page,
		PerPage: size,
	}
	if err := p.Validate(); err != nil {
		return Pagination{}, err
	}
	p.setLink(r, total)
	return p, nil
}

func (p Pagination) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Page, validation.Required, validation.Min(1)),
		validation.Field(&p.PerPage, validation.Required, validation.Min(1)),
	)
}

// SetLinks sets the links of a Pagination object. The fields are conditionally
// set based on the current page value and the total of items.
func (p *Pagination) setLink(r *http.Request, total int) {
	l := Links{}

	if p.Page > 1 {
		l.Prev = buildURLWithPagination(r, *p, -1)
	}
	if total > p.Page*p.PerPage {
		l.Next = buildURLWithPagination(r, *p, 1)
	}

	p.Links = l
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
		"page": p.Page + delta,
	}
	return buildURLWithQuery(r, v)
}
