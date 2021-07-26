package golastic

import (
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Get(ctx ContextConfig, id string) (*esapi.Response, error) {
	res, err := ctx.Client.Get(ctx.IndexName, id)
	if err != nil {
		return nil, ErrUnhandled
	}
	return res, nil
}
