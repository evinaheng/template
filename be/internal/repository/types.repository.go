package repository

// Repository Module
type (
	systemRepo struct {
		ipAddress string
	}
)

type (
	// A System repository provides all queries related for System
	System interface {
		LogOpenFile()
	}
)

// Private structs
