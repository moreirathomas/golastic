package golastic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/clarketm/json" // allows to omit empty structs
)

const (
	defaultOperator  = "and"
	defaultQuerySize = 10
)

// SearchQuery represents an ElasticSearch search query.
// It can be built via NewQuery or NewDefaultQuery.
// It exposes methods to easily retrieve its value
// as bytes, string or via io.Reader.
type SearchQuery struct {
	Query struct {
		MatchAll struct {
			Boost int `json:"boost"`
		} `json:"match_all,omitempty"`
		MultiMatch struct {
			Query    string  `json:"query,omitempty"`
			Operator string  `json:"operator,omitempty"`
			Fields   []Field `json:"fields,omitempty"`
		} `json:"multi_match,omitempty"`
	} `json:"query,omitempty"`
	Sort []map[string]string `json:"sort,omitempty"`
	Size int                 `json:"size,omitempty"`
}

// Bytes returns the raw query as bytes.
func (q SearchQuery) Bytes() []byte {
	// b, _ := json.MarshalIndent(q, "", "  ")
	b, _ := json.Marshal(q)
	return b
}

// String returns the raw query as a string.
func (q SearchQuery) String() string {
	return string(q.Bytes())
}

// Reader returns the raw query as an io.Reader.
func (q SearchQuery) Reader() io.Reader {
	return bytes.NewReader(q.Bytes())
}

// Field is a field name associated with an optional weight.
// It provides marshaling methods allowing to comply automatically
// with ElasticSearch syntax for fields in a query (see MarshalText)
type Field struct {
	Name   string
	Weight int
}

// MarshalText returns the stringified field as a slice of bytes
// and a nil error.
//
// It is automatically called by json.Marshal when it encounters
// a Field value. We use it to format the ElasticSearch query.
//
// For instance, marshaling the following:
//
//	 map[string]interface{}{
//		 "fields": []Field{
//			 {Name: "title", Weight: 10},
//			 {Name: "abstract"}
//		 }
//	 }
//
// results in:
//
// {"fields":["title^10","abstract"]}
func (f Field) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// String returns a string representation of the field in the format
// expected by ElasticSearch.
//
// Examples:
//
// - Field{Name: "title", Weight: 10}.String() == "title^10"
//
// - Field{Name: "abstract"}.String() == "abstract"
func (f Field) String() string {
	if f.Weight == 0 {
		return f.Name
	}
	return fmt.Sprintf("%s^%d", f.Name, f.Weight)
}

// NewDefaultSearchQuery returns a Query targeting all documents
// for the current index, ordered by creation date.
func NewDefaultSearchQuery() SearchQuery {
	q := SearchQuery{}
	q.Query.MatchAll.Boost = 1
	q.Sort = []map[string]string{
		{"_doc": "asc"},
	}
	q.Size = defaultQuerySize

	return q
}

// SearchQueryConfig is a flattened representation of injectable values
// in an ElasticSearch query. The values are then injectected
// in the right place via NewQuery.
// It allows to define a Query conveniently, without having to
// reproduce the whole structure.
type SearchQueryConfig struct {
	Fields []Field
	Sort   []map[string]string
	Size   int
}

// NewSearchQuery returns a Query, built upon the given search query
// and the QueryConfig.
func NewSearchQuery(qs string, cfg SearchQueryConfig) SearchQuery {
	q := SearchQuery{}
	q.Query.MultiMatch.Query = qs
	q.Query.MultiMatch.Fields = cfg.Fields
	q.Query.MultiMatch.Operator = defaultOperator
	q.Size = defaultQuerySize

	if len(cfg.Sort) != 0 {
		q.Sort = cfg.Sort
	}

	if cfg.Size > 0 {
		q.Size = cfg.Size
	}

	return q
}
