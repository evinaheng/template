package usecase

import (
	"github.com/template/be/internal/controller"
)

//  Usecase module
type (
	systemUsecase struct {
		// Controllers
		system controller.System
	}
)

// Usecase interfaces
type (
	// A System usecase provides functions for maintaining system information
	System interface {
		LogOpenFile()
	}
)
