package elasticsearch

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

// Query represents an ElasticSearch search query.
// It can be built via NewQuery or NewDefaultQuery.
// It exposes methods to easily retrieve its value
// as bytes, string or via io.Reader.
type Query struct {
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
func (q Query) Bytes() []byte {
	// b, _ := json.MarshalIndent(q, "", "  ")
	b, _ := json.Marshal(q)
	return b
}

// String returns the raw query as a string.
func (q Query) String() string {
	return string(q.Bytes())
}

// Reader returns the raw query as an io.Reader.
func (q Query) Reader() io.Reader {
	return bytes.NewReader(q.Bytes())
}

// Field is a column name associated with an optional weight.
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

// NewDefaultQuery returns a Query targeting all documents
// for the current index, ordered by creation date.
func NewDefaultQuery() Query {
	q := Query{}
	q.Query.MatchAll.Boost = 1
	q.Sort = []map[string]string{
		{"created_at": "desc"},
		{"_doc": "asc"},
	}
	q.Size = defaultQuerySize

	return q
}

type QueryConfig struct {
	Fields []Field
	Sort   []map[string]string
	Size   int
}

func NewQuery(qs string, cfg QueryConfig) Query {
	q := Query{}
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