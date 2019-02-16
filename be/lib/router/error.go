package router

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/template/be/lib/convert"
)

// NewError for JSON standard API error
func NewError(code string, status int, title string) ArrError {
	err := []Error{}
	err = append(err, Error{
		ID:     code,
		Status: convert.ToString(status),
		Title:  title,
		Code:   status,
	})
	return err
}

// SetError for API result
func (r *HandlerResult) SetError(message, errorCode string, httpStatus int) {
	r.HTTPState = httpStatus
	r.Errors = NewError(errorCode, httpStatus, message)
}

// Error log
func (w ArrError) Error() string {
	errorString, _ := json.Marshal(w)
	return fmt.Sprintf(`{"errors":%s}`, string(errorString))
}

/*
Segregate error codes in array and return an array of error codes & error title
*/
func (w ArrError) Segregate() ([]string, []string) {
	sort.Slice(w, func(i, j int) bool {
		return w[i].ID < w[j].ID
	})

	errorID, errorTitle := []string{}, []string{}

	for _, err := range w {
		errorID = append(errorID, convert.ToString(err.ID))
		errorTitle = append(errorTitle, err.Title)
	}

	return errorID, errorTitle
}
