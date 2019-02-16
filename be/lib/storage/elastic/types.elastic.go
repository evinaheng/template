package elastic

import (
	"context"
	"encoding/json"
)

// An Elastic offers a standard interface for connect to elasticsearch service
type Elastic interface {
	Connect() error
	GetClient() interface{}
	Search(context.Context, SearchParam) ([]json.RawMessage, error)
	Suggest(context.Context, SuggestParam) ([]json.RawMessage, error)
}

// SearchParam for elastic search query
type SearchParam struct {
	Index   string
	Request interface{}
}

// SuggestParam for elastic suggest query
type SuggestParam struct {
	SuggestName string
	Index       string
	Request     interface{}
}
