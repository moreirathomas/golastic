package golastic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/clarketm/json" // allows to omit empty structs
)

// SearchQuery represents an Elasticsearch search query.
// It exposes methods to easily retrieve its value
// as bytes, string or via io.Reader.
type SearchQuery struct {
	// Query represents the body of the full text query being used.
	// Only one of its fields must be used at a time.
	Query struct {
		MatchAll   MatchAllQuery   `json:"match_all,omitempty"`
		MultiMatch MultiMatchQuery `json:"multi_match,omitempty"`
	} `json:"query,omitempty"`

	Sort []map[string]string `json:"sort,omitempty"`
	From int                 `json:"from,omitempty"`
	Size int                 `json:"size,omitempty"`
}

// MatchAllQuery is the query for performing queries
// which matches all documents.
type MatchAllQuery struct {
	Boost int `json:"boost,omitempty"`
}

// MultiMatchQuery is the query for performing full text queries
// accross multiple fields.
type MultiMatchQuery struct {
	Query    string  `json:"query,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
	Operator string  `json:"operator,omitempty"`
}

// Bytes returns the raw query as bytes.
func (q SearchQuery) Bytes() []byte {
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
// with Elasticsearch syntax for fields in a query (see MarshalText)
type Field struct {
	Name   string
	Weight int
}

// MarshalText returns the stringified field as a slice of bytes
// and a nil error.
//
// It is automatically called by json.Marshal when it encounters
// a Field value. We use it to format the Elasticsearch query.
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
// gives:
//
// 	> {"fields":["title^10","abstract"]}
func (f Field) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// String returns a string representation of the field in the format
// expected by Elasticsearch.
// For example:
//
// 	Field{Name: "title", Weight: 10}.String() == "title^10"
// 	Field{Name: "abstract"}.String() == "abstract"
func (f Field) String() string {
	if f.Weight == 0 {
		return f.Name
	}
	return fmt.Sprintf("%s^%d", f.Name, f.Weight)
}

const (
	defaultOperator  = "and"
	defaultQuerySize = 10
)

var (
	defaultSort = []map[string]string{{"_doc": "asc"}}
)

// SearchQueryConfig configures an Elasticsearch full text query.
// Configuration keys are flattened to conveniently define a SearchQuery
// without the need to reproduce its nested structure.
type SearchQueryConfig struct {
	Fields []Field
	Sort   []map[string]string
	From   int // From defines the number of hits to skip.
	Size   int // Size defines the maximum number of hits to return.
}

// MatchAllSearchQuery returns a Query targeting all documents
// for the current index, ordered by creation date.
func MatchAllSearchQuery(size int, from int) SearchQuery {
	q := SearchQuery{}
	// Elasticsearch defaults the boost score to 1 if not provided.
	// q.Query.MatchAll.Boost = 1
	q.Sort = defaultSort
	q.paginate(size, from)
	return q
}

// NewSearchQuery returns a Query, built upon the given search query
// and the QueryConfig.
func NewSearchQuery(qs string, cfg SearchQueryConfig) SearchQuery {
	q := SearchQuery{}

	q.Query.MultiMatch.Query = qs
	q.Query.MultiMatch.Fields = cfg.Fields
	q.Query.MultiMatch.Operator = defaultOperator

	if len(cfg.Sort) != 0 {
		q.Sort = cfg.Sort
	} else {
		q.Sort = defaultSort
	}

	q.paginate(cfg.Size, cfg.From)
	return q
}

func (q *SearchQuery) paginate(size int, from int) {
	if size > 0 {
		q.Size = size
	} else {
		q.Size = defaultQuerySize
	}

	if from >= -1 {
		q.From = from
	}
	// Else use the null value 0
}
