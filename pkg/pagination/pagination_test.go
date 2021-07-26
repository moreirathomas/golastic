package pagination_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moreirathomas/golastic/pkg/pagination"
)

func TestValidate(t *testing.T) {
	// Zero values page and size must error
	_, err := pagination.New(nil, 0, 0, 0)
	if err == nil {
		t.Error("Unexpected nil error: 0 value for parameters \"page\" and \"size\" must error")
	}
}

func TestSetLinks(t *testing.T) {
	// Request page 1, expect page 2 link
	mock := TestLink{
		Request: httptest.NewRequest("GET", buildMockURL(1, 10), nil),
		Prev:    "",
		Next:    buildMockURL(2, 10),
	}
	p, err := pagination.New(mock.Request, 1000, 1, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	mock.makeTest(t, p)

	// Request page 2, expect page 1 and 3 links
	mock = TestLink{
		Request: httptest.NewRequest("GET", buildMockURL(2, 10), nil),
		Prev:    buildMockURL(1, 10),
		Next:    buildMockURL(3, 10),
	}
	p, err = pagination.New(mock.Request, 1000, 2, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	mock.makeTest(t, p)

	// Request page 2, provide a total lesser than the page size, expect page 1 link only
	mock = TestLink{
		Request: httptest.NewRequest("GET", buildMockURL(2, 10), nil),
		Prev:    buildMockURL(1, 10),
		Next:    "",
	}
	p, err = pagination.New(mock.Request, 19, 2, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	mock.makeTest(t, p)
}

type TestLink struct {
	Request *http.Request
	Prev    string
	Next    string
}

func buildMockURL(page, size int) string {
	return fmt.Sprintf("http://localhost:9999/foo?page=%d&size=%d", page, size)
}

func (l TestLink) makeTest(t *testing.T, p pagination.Pagination) {
	if p.Links.Prev != l.Prev {
		t.Fatalf("bad link for prev page: got %s, want %s", p.Links.Prev, l.Prev)
	}
	if p.Links.Next != l.Next {
		t.Fatalf("bad link for next page: got %s, want %s", p.Links.Next, l.Next)
	}
}
