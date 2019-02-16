package dummyelastic

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/template/be/lib/storage/elastic"
)

// New dummy module
func New(config Config) elastic.Elastic {
	return &dummy{
		config: config,
	}
}

// Connect to elastic
func (e *dummy) Connect() error {
	return nil
}

// GetClient of the elastic
func (e *dummy) GetClient() interface{} {
	return nil
}

// Search function
func (e *dummy) Search(ctx context.Context, q elastic.SearchParam) ([]json.RawMessage, error) {

	// Validate empty
	if q.Request == nil {
		return nil, errors.New("Invalid Search parameter")
	}

	request, _ := q.Request.(string)
	resultJSON, ok := e.config.Mocking[request]

	// No result
	if !ok {
		return nil, nil
	}

	results := []json.RawMessage{}
	json.Unmarshal([]byte(resultJSON), &results)
	return results, nil

}

// Suggest function
func (e *dummy) Suggest(ctx context.Context, q elastic.SuggestParam) ([]json.RawMessage, error) {

	// Validate empty
	if q.Request == nil {
		return nil, errors.New("Invalid Suggest parameter")
	}
	request := q.Request.(string)
	resultJSON, ok := e.config.Mocking[request]

	// No result
	if !ok {
		return nil, nil
	}

	results := []json.RawMessage{}
	json.Unmarshal([]byte(resultJSON), &results)
	return results, nil

}
