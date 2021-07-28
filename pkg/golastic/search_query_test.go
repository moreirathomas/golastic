package golastic_test

import (
	"testing"

	"github.com/moreirathomas/golastic/pkg/golastic"
)

func TestMarshaling(t *testing.T) {
	q := golastic.SearchQuery{}
	q.Query.MultiMatch.Fields = []golastic.Field{
		{"title", 10},
		{"abstract", 5},
		{"_doc", 0},
	}

	exp := `{"query":{"multi_match":{"fields":["title^10","abstract^5","_doc"]}}}`

	if got := q.String(); got != exp {
		t.Errorf("unexpected fields marshaling output: expected %s, got %s", exp, got)
	}
}
