package main

import (
	"net/http"
	"time"

	"github.com/template/be/api"
	"github.com/template/be/api/productName/apitest"
	"github.com/template/be/cmd/internal"
	"github.com/template/be/lib/panics"
)

// Initialize all routes for API
func initRoutes(config internal.Config, panicsWrapper panics.Panics, ucase *internal.Usecase) http.Handler {

	// Get environment
	// var isDevelopment = config.Server.Env == internal.Development

	// Middlewares
	// var mw router.Middleware

	// Init ROUTER
	apiRouter := api.NewRouter(api.RouterConfig{
		Timeout: time.Duration(config.Environment.GlobalTimeout) * time.Second,
	})

	// -- Set Panics handler TODO
	apiRouter.PanicsHandler(panicsWrapper.RouterHandler)

	// -----------------------------------------------------------
	// START - ADD ROUTE
	// WARNING! Duplicate routes will cause a panic!
	// -----------------------------------------------------------

	// -- Initialize API Module
	apiTest := apitest.New()

	// -- Routes WITHOUT authentication
	apiRouter.AddRoute(http.MethodGet, "testing/getdata", apiTest.Test)

	// -- Routes WITH authentication

	// -----------------------------------------------------------
	// END - ADD ROUTE
	// -----------------------------------------------------------

	return apiRouter.GetHandler()

}
