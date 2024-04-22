package e_tests

import (
	"elemental/tests/mocks"
	"elemental/tests/setup"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreMeta(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Metadata", t, func() {
		Convey(fmt.Sprintf("Estimated document count should be %d", len(e_mocks.Users)), func() {
			count := UserModel.EstimatedDocumentCount()
			So(count, ShouldEqual, len(e_mocks.Users))
		})
	})
}
