package golastic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/clarketm/json" // allows to omit empty structs
)

const (
	defaultOperator  = "and"
	DefaultQuerySize = 10
	DefaultQueryFrom = 0
)

var defaultSort = []string{"_doc:asc"}

// SearchPagination configures the pagination of an Elasticsearch search query.
type SearchPagination struct {
	From int
	Size int
}

// SearchQuery configures the sort parameter of an Elasticsearch search query.
type SearchSort []string

// SearchQuery represents the body of query made to Elasticsearch
// Search API. It is shaped as expected from Elasticsearch.
//
// It exposes methods for easy conversion to bytes, string or io.Reader.
//
// The nested field Query holds the full text query being used.
// Only one of its fields must be used at a time.
type SearchQuery struct {
	Query struct {
		MatchAll   MatchAllQuery   `json:"match_all,omitempty"`
		MultiMatch MultiMatchQuery `json:"multi_match,omitempty"`
	} `json:"query,omitempty"`
}

// MatchAllQuery is the query for performing queries
// which match all documents.
type MatchAllQuery struct {
	Boost int `json:"boost,omitempty"`
}

// MultiMatchQuery is the query for performing full text queries
// across multiple fields.
type MultiMatchQuery struct {
	Query    string  `json:"query,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
	Operator string  `json:"operator,omitempty"`
}

// Bytes returns the query as bytes.
func (q SearchQuery) Bytes() []byte {
	b, _ := json.Marshal(q)
	return b
}

// String returns the query as a string.
func (q SearchQuery) String() string {
	return string(q.Bytes())
}

// Reader returns the query as an io.Reader.
func (q SearchQuery) Reader() io.Reader {
	return bytes.NewReader(q.Bytes())
}

// Field represents a field as expected by Elasticsearch for
// multi-match queries. A field can optionally have a weight.
//
// It provides marshaling methods allowing to comply automatically
// with Elasticsearch syntax (see Field.MarshalText).
type Field struct {
	Name   string
	Weight int
}

// MarshalText returns a field stringified as a slice of bytes
// and a nil error.
//
// It is automatically called by json.Marshal when it encounters
// a Field value. It is used when formatting an Elasticsearch query.
//
// For instance, marshaling the following:
//
//	map[string]interface{}{
//		"fields": []Field{
//			{Name: "title", Weight: 10},
//			{Name: "abstract"}
//		}
//	}
//
// gives:
//
//	{"fields":["title^10","abstract"]}
func (f Field) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// String returns a field as a string formatted as expected by Elasticsearch.
// For example:
//
//	Field{Name: "title", Weight: 10}.String() == "title^10"
//	Field{Name: "abstract"}.String() == "abstract"
func (f Field) String() string {
	if f.Weight == 0 {
		return f.Name
	}
	return fmt.Sprintf("%s^%d", f.Name, f.Weight)
}

// newMatchAllQuery returns a configured SearchQuery for match-all queries.
func newMatchAllQuery() SearchQuery {
	q := SearchQuery{}
	q.Query.MatchAll.Boost = 1
	return q
}

// newMultiMatchQuery returns a configured SearchQuery for multi-match queries.
func newMultiMatchQuery(qs string, f []Field) SearchQuery {
	q := SearchQuery{}
	q.Query.MultiMatch.Query = qs
	q.Query.MultiMatch.Fields = f
	q.Query.MultiMatch.Operator = defaultOperator
	return q
}
