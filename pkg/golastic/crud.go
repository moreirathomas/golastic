package golastic

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Get(ctx ContextConfig, id string) (*esapi.Response, error) {
	res, err := ctx.Client.Get(ctx.IndexName, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}
	return res, nil
}

func Update(ctx ContextConfig, id string, doc interface{}) (*esapi.Response, error) {
	// document must be wrapped in a "doc" object
	payload, err := json.Marshal(map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := ctx.Client.Update(ctx.IndexName, id, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

func Insert(ctx ContextConfig, doc interface{}) (*esapi.Response, error) {
	payload, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := ctx.Client.Index(ctx.IndexName, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

func Delete(ctx ContextConfig, id string) (*esapi.Response, error) {
	res, err := ctx.Client.Delete(ctx.IndexName, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}
