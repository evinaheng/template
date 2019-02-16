package router

import (
	"context"
	"net/http"
)

type (

	// HTTPHandler return function for API call
	HTTPHandler func(w http.ResponseWriter, r *http.Request) error

	// A Router provides interface for HTTP routing function
	Router interface {

		// Add new route
		// Param : method, string, function call
		// Function call param : version
		AddRoute(string, string, func(int) func(context.Context, *HandlerParam) HandlerResult)

		// Use middlware for all routes
		Use(Middleware)

		// Get router handler
		GetHandler() http.Handler

		// With use inline middleware to specific routes
		With(...Middleware) Router

		// PanicHandler for handling panics
		PanicsHandler(func(context.Context, *http.Request, *HandlerParam, func(context.Context, *HandlerParam) HandlerResult) HandlerResult)
	}

	// Handler for API routing
	Handler interface {
		Do(context.Context, *HandlerParam) HandlerResult
	}

	// Middleware interface
	Middleware interface {
		Handler(http.Handler) http.Handler
	}

	// HandlerParam for Handler
	HandlerParam struct {
		Lang string            // Language
		Vars map[string]string // Custom variables
	}

	// HandlerResult for Handler
	HandlerResult struct {
		JSON         *Data             // Response JSON - High Priority
		Print        []byte            // Response print
		HTTPState    int               // Response HTTP status
		Errors       ArrError          // Response result for errors
		Headers      map[string]string // Custom headers for response
		contextError error             // For context error
	}

	// Error struct for failed API process
	Error struct {
		ID     string      `json:"id,omitempty"`
		Status string      `json:"status,omitempty"`
		Title  string      `json:"title,omitempty"`
		Code   int         `json:"-"`
		Meta   interface{} `json:"meta,omitempty"`
	}

	// ArrError array of Error
	ArrError []Error

	// Data structure for API
	Data struct {
		Data     interface{} `json:"data,omitempty"`
		Included interface{} `json:"included,omitempty"`
		Errors   ArrError    `json:"errors,omitempty"`
		Meta     interface{} `json:"meta,omitempty"`
	}

	// Object structure inside data
	Object struct {
		Type          string      `json:"type,omitempty"`
		ID            string      `json:"id,omitempty"`
		Attributes    interface{} `json:"attributes,omitempty"`
		Relationships interface{} `json:"relationships,omitempty"`
	}
	// ArrObject array of objects for sort implementation
	ArrObject []Object

	// Relationship structure for API
	Relationship struct {
		Data RelationshipData `json:"data"`
	}

	// Relationships array of Relationship
	Relationships struct {
		Data []RelationshipData `json:"data"`
	}

	// RelationshipData data structure inside relationship
	RelationshipData struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}
)

func (a ArrObject) Len() int           { return len(a) }
func (a ArrObject) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArrObject) Less(i, j int) bool { return a[i].ID < a[j].ID }
