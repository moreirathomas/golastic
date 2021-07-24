package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moreirathomas/golastic/pkg/httputil"
)

func TestNewPagination(t *testing.T) {
	mock := []struct {
		httputil.Pagination
		expected int
	}{
		{httputil.NewPagination(10, 0), 1},
		{httputil.NewPagination(10, 10), 2},
		{httputil.NewPagination(10, 11), 2},
		{httputil.NewPagination(10, 100), 11},
	}

	for _, v := range mock {
		if v.Page != v.expected {
			t.Fatalf("unexpected page value: got %d, want %d", v.Page, v.expected)
		}
	}
}

func TestSetLinks(t *testing.T) {
	mockRequest := func(target string) *http.Request {
		return httptest.NewRequest("GET", target, nil)
	}

	type expect struct {
		prev string
		next string
	}

	mock := []struct {
		request    *http.Request
		pagination httputil.Pagination
		expected   expect
	}{
		{
			request:    mockRequest("http://localhost:9999/books?query=foo&size=1&from=0"),
			pagination: httputil.NewPagination(1, 0),
			expected: expect{
				prev: "",
				next: "http://localhost:9999/books?from=1&query=foo&size=1"},
		},
		{
			request:    mockRequest("http://localhost:9999/books?query=foo&size=1&from=1"),
			pagination: httputil.NewPagination(1, 1),
			expected: expect{
				prev: "http://localhost:9999/books?from=0&query=foo&size=1",
				next: "http://localhost:9999/books?from=2&query=foo&size=1"},
		},
		{
			request:    mockRequest("http://localhost:9999/books?query=foo&size=10&from=10"),
			pagination: httputil.NewPagination(10, 10),
			expected: expect{
				prev: "http://localhost:9999/books?from=0&query=foo&size=10",
				next: "http://localhost:9999/books?from=20&query=foo&size=10"},
		},
	}

	for _, v := range mock {
		// set total to 1000 to ensure we always have a next field
		v.pagination.SetLinks(v.request, 1000)

		if v.pagination.Link.Prev != v.expected.prev {
			t.Fatalf("bad link for prev page: got %s, want %s", v.pagination.Link.Prev, v.expected.prev)
		}
		if v.pagination.Link.Next != v.expected.next {
			t.Fatalf("bad link for next page: got %s, want %s", v.pagination.Link.Next, v.expected.next)
		}
	}
}
