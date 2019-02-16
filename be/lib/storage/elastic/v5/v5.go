package v5

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/template/be/lib/storage/elastic"
	eslib "gopkg.in/olivere/elastic.v5"
)

// New v5 elastic module
func New(config Config) elastic.Elastic {
	return &elasticV5{
		config: config,
	}
}

// Connect to elastic
func (e *elasticV5) Connect() error {

	var err error

	// Connect to elastic client
	e.client, err = eslib.NewClient(
		eslib.SetURL(e.config.Endpoint),
		eslib.SetSniff(false),
		eslib.SetHealthcheck(false))

	return err
}

// GetClient of the elastic
func (e *elasticV5) GetClient() interface{} {
	return e.client
}

// Search function
func (e *elasticV5) Search(ctx context.Context, q elastic.SearchParam) ([]json.RawMessage, error) {

	// Validate empty
	if q.Index == "" || q.Request == nil {
		return nil, errors.New("Invalid Search parameter")
	}

	searchService := e.client.Search().Index(q.Index).Type(q.Index)
	queryResult, err := searchService.Source(q.Request).Do(ctx)
	if err != nil {
		return nil, err
	}

	results := []json.RawMessage{}
	for _, h := range queryResult.Hits.Hits {
		results = append(results, *h.Source)
	}

	return results, nil

}

// Suggest function
func (e *elasticV5) Suggest(ctx context.Context, q elastic.SuggestParam) ([]json.RawMessage, error) {

	// Validate empty
	if q.Index == "" || q.Request == nil || q.SuggestName == "" {
		return nil, errors.New("Invalid Suggest parameter")
	}

	searchService := e.client.Search().Index(q.Index).Type(q.Index)
	queryResult, err := searchService.Source(q.Request).Do(ctx)
	if err != nil {
		return nil, err
	}

	results := []json.RawMessage{}

	// Process elastic suggestions
	if suggestions, ok := queryResult.Suggest[q.SuggestName]; ok {

		for r := range suggestions {

			for o := range suggestions[r].Options {
				results = append(results, *suggestions[r].Options[o].Source)

			}

		}

	}

	return results, nil

}
