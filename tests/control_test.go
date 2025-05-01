package e_tests

import (
	"os"
	"testing"

	e_test_setup "github.com/elcengine/elemental/tests/setup"
)

func TestMain(m *testing.M) {
    code := m.Run()

    e_test_setup.Teardown()

    os.Exit(code)
}