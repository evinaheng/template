package api

import (
	"net/http"
)

// Available API versions
var (
	versionsString = []string{"v1", "v1.1"}
	versions       = []int{1, 2}
)

/*
writeResponseJSON for result of API
*/
func writeResponseJSON(w http.ResponseWriter, data []byte, httpState int, headers map[string]string) {
	// Set HTTP Status if isn't 0
	// If it's zero, app will set 200
	if httpState > 0 {
		w.WriteHeader(httpState)
	}

	if len(headers) == 0 {

		// Set header application/json
		w.Header().Set("Content-Type", "application/json")

	} else {

		// Set custom headers
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.Write(data)
}
