package internal_test

import (
	"testing"

	"github.com/moreirathomas/golastic/internal"
)

func TestNewPagination(t *testing.T) {
	mock := []struct {
		internal.Pagination
		want int
	}{
		{internal.NewPagination(10, 0), 1},
		{internal.NewPagination(10, 10), 2},
		{internal.NewPagination(10, 11), 2},
		{internal.NewPagination(10, 100), 11},
	}

	for _, v := range mock {
		if v.Page != v.want {
			t.Fatalf("unexpected page value: got %d, want %d", v.Page, v.want)
		}
	}
}
