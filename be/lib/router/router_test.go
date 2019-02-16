package router_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/router"
)

func TestServeHTTP(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "https://site.com", nil)
	w := httptest.NewRecorder()

	result := HTTPHandler(func(w http.ResponseWriter, r *http.Request) error {
		return NewError("0", http.StatusInternalServerError, "Error here")
	})
	result.ServeHTTP(w, r)
	r.Method = http.MethodOptions
	result.ServeHTTP(w, r)
}

func TestServeHTTPError(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "https://site.com", nil)
	w := httptest.NewRecorder()

	result := HTTPHandler(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("Test")
	})
	result.ServeHTTP(w, r)
	r.Method = http.MethodOptions
	result.ServeHTTP(w, r)
}

func TestSortObjects(t *testing.T) {

	arr := ArrObject{}
	arr = append(arr, Object{
		ID: "2",
	})
	arr = append(arr, Object{
		ID: "1",
	})

	sort.Sort(arr)

	assert.Equal(t, "1", arr[0].ID)
	assert.Equal(t, "2", arr[1].ID)
}

func TestSetContextError(t *testing.T) {
	hr := HandlerResult{}
	hr.SetContextError(errors.New("foo"))
	assert.EqualError(t, hr.GetContextError(), "foo")
}
