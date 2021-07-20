package golastic_test

import (
	"testing"

	"github.com/moreirathomas/golastic/pkg/golastic"
)

func TestMarshaling(t *testing.T) {
	test := marshalingTest{}
	t.Run("Fields marshaling", test.fields)
	t.Run("Sort marshaling", test.sort)
}

type marshalingTest struct{}

func (test marshalingTest) fields(t *testing.T) {
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

func (test marshalingTest) sort(t *testing.T) {
	q := golastic.SearchQuery{}
	q.Sort = []map[string]string{
		{"_score": "desc"},
		{"_doc": "asc"},
		{"created_at": "desc"},
	}

	exp := `{"sort":[{"_score":"desc"},{"_doc":"asc"},{"created_at":"desc"}]}`

	if got := q.String(); got != exp {
		t.Errorf("unexpected sort marshaling output: expected %s, got %s", exp, got)
	}
}
