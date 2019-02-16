package router_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/router"
)

func TestError(t *testing.T) {
	var a ArrError
	assert.Equal(t, `{"errors":null}`, a.Error())
}

func TestNewError(t *testing.T) {
	apiError := NewError("1", 500, "Testing Error")
	assert.Equal(t, `{"errors":[{"id":"1","status":"500","title":"Testing Error"}]}`, apiError.Error())

}

func TestSegregate(t *testing.T) {
	var errors ArrError

	errorA := Error{
		ID:     "1",
		Status: "400",
		Title:  "Testing Error A",
		Code:   400,
	}

	errorB := Error{
		ID:     "2",
		Status: "500",
		Title:  "Testing Error B",
		Code:   500,
	}

	errors = append(errors, errorA)
	errors = append(errors, errorB)

	errorID, errorTitles := errors.Segregate()
	assert.Equal(t, []string{errorA.ID, errorB.ID}, errorID)
	assert.Equal(t, []string{errorA.Title, errorB.Title}, errorTitles)
}

func TestSetError(t *testing.T) {

	result := &HandlerResult{}
	result.SetError("foo", "123", http.StatusInternalServerError)

	assert.Equal(t, "foo", result.Errors[0].Title)
	assert.Equal(t, "123", result.Errors[0].ID)
	assert.Equal(t, http.StatusInternalServerError, result.Errors[0].Code)
	assert.Equal(t, http.StatusInternalServerError, result.HTTPState)

}
