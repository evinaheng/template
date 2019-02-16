package v5_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/template/be/lib/storage/elastic"
	. "github.com/template/be/lib/storage/elastic/v5"
)

func TestConnectError(t *testing.T) {

	config := Config{}
	es := New(config)

	assert.EqualError(t, es.Connect(), "no active connection found: no Elasticsearch node available")
	assert.Nil(t, es.GetClient())
}

func TestConnectSuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"name" : "elastic",
			"cluster_name" : "elastic",
			"cluster_uuid" : "XqzEn1WWQMa6dGZmdS3A1g",
			"version" : {
			  "number" : "5.2.2",
			  "build_hash" : "f9d9b74",
			  "build_date" : "2017-02-24T17:26:45.835Z",
			  "build_snapshot" : false,
			  "lucene_version" : "6.4.1"
			},
			"tagline" : "You Know, for Search"
		  }
		  `
		w.Write([]byte(resp))
	}

	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)

	assert.Nil(t, es.Connect())

	assert.NotNil(t, es.GetClient())
}

func TestSearchQueryError(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `boom`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Invalid param
	param := elastic.SearchParam{}
	result, err := es.Search(context.TODO(), param)
	assert.EqualError(t, err, "Invalid Search parameter")

	// Error during hit elastic
	param = elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	result, err = es.Search(context.TODO(), param)
	assert.Empty(t, result)
	assert.EqualError(t, err, "invalid character 'b' looking for beginning of value")
}

func TestSearchQueryEmptySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 4,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": null,
					"hits": []
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	result, err := es.Search(context.TODO(), param)
	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestSearchQuerySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 1,
				"timed_out": false,
				"_shards": {
					"total": 5,"successful": 5,"skipped": 0,"failed": 0
				},
				"hits": {
					"total": 2,"max_score": 1,
					"hits": [
						{
							"_index": "halo",
							"_type": "halo",
							"_id": "35b4f0710c0d0e1b3fda664a27f27f38",
							"_score": 1,
							"_source": {
								"key": "key0",
								"agent_id": 0					
							}
						},
						{
							"_index": "halo",
							"_type": "halo",
							"_id": "c8e0f4477f8c2884010e697dd460a70c",
							"_score": 1,
							"_source": {
								"key": "key1",
								"agent_id": 1
							}
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SearchParam{
		Index:   "halo",
		Request: "foo",
	}
	result, err := es.Search(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	for k, v := range result {
		res := dummyResult{}
		json.Unmarshal(v, &res)

		assert.Equal(t, k, res.AgentID)
		assert.Equal(t, fmt.Sprintf("key%d", k), res.Key)

	}
}

func TestSuggestQueryError(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `boom`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Invalid param
	param := elastic.SuggestParam{}
	result, err := es.Suggest(context.TODO(), param)
	assert.EqualError(t, err, "Invalid Suggest parameter")

	// Error during hit elastic
	param = elastic.SuggestParam{
		Index:       "halo",
		Request:     "foo",
		SuggestName: "bar",
	}
	result, err = es.Suggest(context.TODO(), param)
	assert.Empty(t, result)
	assert.EqualError(t, err, "invalid character 'b' looking for beginning of value")
}

func TestSuggestQueryEmptySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 4,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"search-suggest": [
						{
							"text": "xxxxxxxxxxxxxxxx",
							"offset": 0,
							"length": 9,
							"options": []
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SuggestParam{
		SuggestName: "search-suggest",
		Index:       "halo",
		Request:     "foo",
	}
	result, err := es.Suggest(context.TODO(), param)
	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestSuggestQuerySuccess(t *testing.T) {

	// Mocking HTTP
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {

		resp := `{"name":"elastic","cluster_name":"elastic","cluster_uuid":"XqzEn1WWQMa6dGZmdS3A1g","version":{"number":"5.2.2","build_hash":"f9d9b74","build_date":"2017-02-24T17:26:45.835Z","build_snapshot":false,"lucene_version":"6.4.1"},"tagline":"You Know, for Search"}`

		if r.URL.String() == "/halo/halo/_search" {
			resp = `{
				"took": 150,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"search-suggest": [
						{
							"text": "agent",
							"offset": 0,
							"length": 3,
							"options": [
								{
									"text": "Agent1",
									"_index": "halo",
									"_type": "halo",
									"_id": "2",
									"_score": 152,
									"_source": {
										"key": "key0",
										"agent_id": 0
									}
								}
							]
						},
						{
							"text": "agent",
							"offset": 0,
							"length": 3,
							"options": [
								{
									"text": "Agent1",
									"_index": "halo",
									"_type": "halo",
									"_id": "3",
									"_score": 152,
									"_source": {
										"key": "key1",
										"agent_id": 1
									}
								}
							]
						}
					]
				}
			}`
		}

		w.Write([]byte(resp))
	}

	// Connect to elastic
	config := Config{
		Endpoint: ts.URL,
	}
	es := New(config)
	es.Connect()

	// Error during hit elastic
	param := elastic.SuggestParam{
		SuggestName: "search-suggest",
		Index:       "halo",
		Request:     "foo",
	}
	result, err := es.Suggest(context.TODO(), param)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
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
