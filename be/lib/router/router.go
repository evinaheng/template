package router

import (
	"fmt"
	"net/http"
)

/*
ServeHTTP listen and serve request
*/
func (fn HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	if err := fn(w, r); err != nil {
		sendError(w, r, err)
	}

}

/*
sendError print response to browser
*/
func sendError(w http.ResponseWriter, r *http.Request, err error) {

	var wErr []Error
	var wErrObj ArrError
	var ok bool

	// Standard API Error
	if wErr, ok = err.(ArrError); ok && len(wErr) > 0 {
		wErrObj = wErr
	} else {
		wErrObj = NewError("0", http.StatusInternalServerError, "Unknown Error")
		wErr = wErrObj
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(wErr[0].Code)
	fmt.Fprint(w, wErrObj.Error())

}

// SetContextError to handler
func (h *HandlerResult) SetContextError(err error) {
	h.contextError = err
}

// GetContextError to handler
func (h *HandlerResult) GetContextError() error {
	return h.contextError
}
