//go:build benchmark

package benchmarks

import (
	"os"
	"testing"

	ts "github.com/elcengine/elemental/tests/fixtures/setup"
)

func TestMain(m *testing.M) {
	code := m.Run()

	ts.Teardown()

	os.Exit(code)
}
