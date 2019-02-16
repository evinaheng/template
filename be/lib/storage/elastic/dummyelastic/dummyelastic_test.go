package dummyelastic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/template/be/lib/storage/elastic"
	. "github.com/template/be/lib/storage/elastic/dummyelastic"
)

func TestGetClientSuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	assert.Nil(t, es.Connect())
	assert.Nil(t, es.GetClient())
}

func TestSearchQueryError(t *testing.T) {

	config := Config{}
	es := New(config)

	// Invalid param
	param := elastic.SearchParam{}
	result, err := es.Search(context.TODO(), param)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Invalid Search parameter")

}

func TestSearchQueryEmptySuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	// Without mocking
	param := elastic.SearchParam{
		Request: "empty response",
	}
	result, err := es.Search(context.TODO(), param)
	assert.Empty(t, result)
	assert.Nil(t, err)

}

func TestSearchQuerySuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = `
	[{
		"key": "key0",
		"agent_id": 0					
	},
	{
		"key": "key1",
		"agent_id": 1
	}]`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.SearchParam{
		Request: "foo",
	}
	result, err := es.Search(context.TODO(), param)

	assert.Nil(t, err)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)

	}

}

func TestSuggestQueryError(t *testing.T) {

	config := Config{}
	es := New(config)

	// Invalid param
	param := elastic.SuggestParam{}
	result, err := es.Suggest(context.TODO(), param)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Invalid Suggest parameter")

}

func TestSuggestQueryEmptySuccess(t *testing.T) {

	config := Config{}
	es := New(config)

	// Without mocking
	param := elastic.SuggestParam{
		Request: "empty response",
	}
	result, err := es.Suggest(context.TODO(), param)
	assert.Empty(t, result)
	assert.Nil(t, err)

}

func TestSuggestQuerySuccess(t *testing.T) {

	mocking := map[string]string{}

	mocking["foo"] = `
	[{
		"key": "key0",
		"agent_id": 0					
	},
	{
		"key": "key1",
		"agent_id": 1
	}]`
	config := Config{
		Mocking: mocking,
	}
	es := New(config)

	// With mocking
	param := elastic.SuggestParam{
		Request: "foo",
	}
	result, err := es.Suggest(context.TODO(), param)

	assert.Nil(t, err)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)

	}

}

type dummyResult struct {
	Key     string `json:"key"`
	AgentID int    `json:"agent_id"`
}
