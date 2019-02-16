package apitest

import (
	"context"

	"github.com/template/be/lib/router"
)

var instance *Module

// Module for APIDropdown
type Module struct {
	// put all usecase used
}

// New module
func New() *Module {
	return &Module{}
}

// Define ALL Related Module for Api Test here

// API Module
func (m *Module) Test(version int) func(ctx context.Context, param *router.HandlerParam) (apiresult router.HandlerResult) {
	return m.getTestingData
}
