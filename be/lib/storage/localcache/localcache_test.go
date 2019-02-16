package localcache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/storage/localcache"
)

func TestLocalcache(t *testing.T) {
	// Set value
	Set("a", 1, 11)

	// Get value
	result, ok := Get("a")
	assert.Equal(t, 1, result)
	assert.True(t, ok)

	// Delete value
	Delete("a")
}
