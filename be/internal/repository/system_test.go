package repository_test

import (
	"testing"

	. "github.com/template/be/internal/repository"
)

func TestSystemLogOpenFile(t *testing.T) {
	systemRepo := NewSystem("")
	systemRepo.LogOpenFile()
}
