package locale

import (
	"os"
	"testing"

	"github.com/template/be/lib/constants/env"
)

func TestMain(m *testing.M) {

	Init(env.Development)
	os.Exit(m.Run())
}
