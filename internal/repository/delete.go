package repository

import (
	"fmt"
)

func (r *Repository) Delete(id string) error {
	res, err := r.es.Delete(r.indexName, id)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}
