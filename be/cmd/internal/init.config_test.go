package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	ac := &Config{}
	ok := readConfig(ac, "path", "test")
	assert.False(t, ok)
}
