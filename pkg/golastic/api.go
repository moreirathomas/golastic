package golastic

import "github.com/elastic/go-elasticsearch/v7"

// ContextConfig configures the context for a Elasticsearch API call.
type ContextConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}

// Indices interfaces Elasticsearch Indices API.
func Indices(c *elasticsearch.Client) IndicesAPI {
	return IndicesAPI{client: c}
}

// Document interfaces Elasticsearch Document API.
func Document(cfg ContextConfig) DocumentAPI {
	return DocumentAPI{ContextConfig{
		Client:    cfg.Client,
		IndexName: cfg.IndexName,
	}}
}

// Search interfaces Elasticsearch Search API.
func Search(cfg ContextConfig) SearchAPI {
	return SearchAPI{ContextConfig{
		Client:    cfg.Client,
		IndexName: cfg.IndexName,
	}}
}
