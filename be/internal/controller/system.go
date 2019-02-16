package controller

import "github.com/template/be/internal/repository"

// NewSystem New system service
func NewSystem(systemRP repository.System) System {

	svc := &systemSvc{
		systemRP: systemRP,
	}

	return svc
}

func (s *systemSvc) LogOpenFile() {
	s.systemRP.LogOpenFile()
}
