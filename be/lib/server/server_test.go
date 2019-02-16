package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/server"
)

func TestGetIPAddress(t *testing.T) {
	assert.NotNil(t, GetIPAddress())
}
