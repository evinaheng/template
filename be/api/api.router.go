/*Package api router
 */
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/template/be/ctxs"
	"github.com/template/be/lib/language"
	"github.com/template/be/lib/router"
	"github.com/template/be/locale"
)

var (
	regexPathname = regexp.MustCompile(`(/:[\w_]{1,})`)
)

// NewRouter for API
func NewRouter(config RouterConfig) router.Router {

	// Set timeout - default is 29s
	if config.Timeout == time.Duration(0) {
		config.Timeout = time.Duration(29 * time.Second)
	}

	// Return  router
	return &apiRouter{
		rt:     httprouter.New(),
		config: config,
	}

}

// With use inline middleware to specific routes
func (fr *apiRouter) With(inlineMw ...router.Middleware) router.Router {

	// Remove nil middleware
	var mwArr = []router.Middleware{}
	for i := range inlineMw {
		if inlineMw[i] != nil {
			mwArr = append(mwArr, inlineMw[i])
		}
	}

	// Return router
	return &apiRouter{
		rt:                fr.rt,
		config:            fr.config,
		middlewares:       fr.middlewares,
		panicsHandler:     fr.panicsHandler,
		inlineMiddlewares: mwArr,
	}
}

// AddRoute add new path to router
func (fr *apiRouter) AddRoute(method, path string, f func(int) func(context.Context, *router.HandlerParam) router.HandlerResult) {

	// Generate pathname for tracking purpose
	pathname := generatePathname(path)

	/*
		Generate route handling per API version
		e.g. v1/search/single
		e.g. v1.1/search/single
	*/
	for index, version := range versionsString {

		// Get new handler for API
		newHandler := router.HTTPHandler(fr.createAPIHandler(versions[index], pathname, method, f))

		var wrappedHandler http.Handler
		wrappedHandler = newHandler

		// Inject middlewares
		for i := range fr.middlewares {
			wrappedHandler = fr.middlewares[i].Handler(wrappedHandler)
		}

		// Inject inline middlewares
		for i := range fr.inlineMiddlewares {
			wrappedHandler = fr.inlineMiddlewares[i].Handler(wrappedHandler)
		}

		withParamHandler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

			/*
				Initialize variables in URL
				e.g. cart/{id} -> cart/123 -> vars[id] is 123
			*/
			vars := map[string]string{}
			for i := range ps {
				vars[ps[i].Key] = ps[i].Value
			}

			// Context is cancelled by user before we start func call -- exit immediately
			if r.Context().Err() != nil && r.Context().Err().Error() == ctxs.ContextCancelledError {
				return
			}

			// Inject vars to context
			ctx := context.WithValue(r.Context(), ctxs.RouterVarContextKey, vars)

			wrappedHandler.ServeHTTP(w, r.WithContext(ctx))
		}

		fr.rt.Handle(method, "/"+version+"/"+path+"/", withParamHandler)

	}

}

// Use middleware for all routes
func (fr *apiRouter) Use(mw router.Middleware) {
	fr.middlewares = append(fr.middlewares, mw)
}

// PanicsHandler set middleware for panics
func (fr *apiRouter) PanicsHandler(f func(context.Context, *http.Request, *router.HandlerParam, func(context.Context, *router.HandlerParam) router.HandlerResult) router.HandlerResult) {
	fr.panicsHandler = f
}

// GetHandler return http handler
func (fr *apiRouter) GetHandler() http.Handler {
	return fr.rt
}

// Get pathname for tracking
func generatePathname(url string) string {
	pathname := regexPathname.ReplaceAllString(url, "")
	pathname = strings.Replace(pathname, "//", "/", -1)
	pathname = strings.Replace(pathname, "//", "/", -1)
	pathname = strings.Trim(pathname, "/")
	pathname = strings.Replace(pathname, "/", ".", -1)
	return pathname
}

// create new handler for API
func (fr *apiRouter) createAPIHandler(version int, pathname, method string, f func(int) func(context.Context, *router.HandlerParam) router.HandlerResult) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {

		// Set all requests to context + initialize OpenTracing
		var operationName = "api." + pathname + "." + method
		span, ctx := opentracing.StartSpanFromContext(ctxs.GetAllContextFromRequest(r), operationName)
		defer span.Finish()

		// Get router variables
		// e.g. test/{id} -> test/123 -> vars[id] is 123
		vars, _ := ctx.Value(ctxs.RouterVarContextKey).(map[string]string)

		// Get Language parameter
		currLang, _ := ctxs.LanguageHeaderKeyFromContext(ctx)
		currLang = language.GetDefault(currLang)

		// Prepare handler paramater
		doParam := &router.HandlerParam{
			Lang: currLang,
			Vars: vars,
		}

		// Call API function
		result := fr.callWithTimeout(ctx, r, version, doParam, f)

		// Response for client
		var responseData []byte
		var responseJSON interface{}

		// If there's context error
		if result.GetContextError() != nil {

			if result.GetContextError().Error() == ctxs.ContextCancelledError {
				// Cancelled by user
				responseJSON = router.Data{
					Errors: router.NewError("97", http.StatusBadRequest, locale.Translate(currLang, "ErrorServer")),
				}
				result.HTTPState = 499 // Nginx error for context cancelled

			} else {
				// Timeout
				responseJSON = router.Data{
					Errors: router.NewError("14", http.StatusGatewayTimeout, locale.Translate(currLang, "ErrorServer")),
				}
				result.HTTPState = http.StatusGatewayTimeout // Set 504 error
			}

		} else if result.JSON != nil {
			// High Priority for JSON response
			// JSON Data exists -- No Error / Happy Flow
			responseJSON = result.JSON

		} else if result.Print != nil {
			// Low Priority for Print response
			// Print Data exists -- No Error / Happy Flow
			responseData = result.Print

		} else if result.HTTPState != 0 {
			// Error exists
			responseJSON = router.Data{
				Errors: result.Errors,
			}

		} else {

			// Panics
			responseJSON = router.Data{
				Errors: router.NewError("99", http.StatusInternalServerError, locale.Translate(currLang, "ErrorServer")),
			}
			result.HTTPState = http.StatusInternalServerError // Set 500 error

		}

		// Use ResponseJSON for response
		if responseData == nil && responseJSON != nil {
			responseData, _ = json.Marshal(responseJSON)
		}

		// Write response to client
		writeResponseJSON(w, responseData, result.HTTPState, result.Headers)
		return nil
	}
}

// callWithTimeout call API function with timeout
func (fr *apiRouter) callWithTimeout(ctx context.Context, r *http.Request, version int, doParam *router.HandlerParam, f func(int) func(context.Context, *router.HandlerParam) router.HandlerResult) router.HandlerResult {

	// Add timeout to context
	ct, cancel := context.WithTimeout(ctx, fr.config.Timeout)
	defer cancel()

	// Use channel for context
	ch := make(chan router.HandlerResult, 1)
	var result router.HandlerResult

	go func() {

		// Use panic wrapper
		ch <- fr.panicsHandler(ctx, r, doParam, f(version))

	}()

	select {
	case <-ct.Done():
		// Timeout or cancelled
		res := router.HandlerResult{}
		res.SetContextError(ct.Err())
		return res

	case result = <-ch:
	}

	return result
}
