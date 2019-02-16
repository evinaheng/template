package apitest

import (
	"os"
	"testing"
)

var testModule *Module

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
