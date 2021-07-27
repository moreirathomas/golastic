package golastic

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// esInsertResponse represents the structure of an Elasticsearch response
// for an insertion request.
type esInsertResponse struct {
	Result string `json:"result"` // "created" in case of success
	ID     string `json:"_id"`    // the document generated ID
}

// unwrap returns the document ID wrapped in esGetResponse
// or the first non-nil error encountered in the process.
func (r esInsertResponse) unwrap() (string, error) {
	if r.Result != "created" {
		// TODO: error handling
		return "", errors.New("not created")
	}

	return r.ID, nil
}

// ReadInsertResponse reads an Elasticsearch response for an insertion request
// and returns the document ID or the first non-nil error occurring in the process.
func ReadInsertResponse(res *esapi.Response) (string, error) {
	if err := ReadErrorResponse(res); err != nil {
		return "", err
	}

	r, err := decodeRawInsertResponse(res)
	if err != nil {
		return "", ErrUnhandled
	}

	return r.unwrap()
}

func decodeRawInsertResponse(res *esapi.Response) (esInsertResponse, error) {
	var r esInsertResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}
