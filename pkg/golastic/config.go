package golastic

import "github.com/elastic/go-elasticsearch/v7"

// ContextConfig configures the context for a Elasticsearch API call.
type ContextConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}
