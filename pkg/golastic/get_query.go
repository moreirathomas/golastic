package golastic

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Get(client *elasticsearch.Client, indexName, id string) (*esapi.Response, error) {
	res, err := client.Get(indexName, id)
	if err != nil {
		return nil, ErrUnhandled
	}
	return res, nil
}
