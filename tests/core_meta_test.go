package e_tests

import (
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreMeta(t *testing.T) {
	t.Parallel()

	var LocalUserModel = UserModel.Clone().SetCollection("users_for_meta")

	e_test_setup.Connection()

	LocalUserModel.InsertMany(e_mocks.Users).Exec()

	defer e_test_setup.Teardown()

	Convey("Metadata", t, func() {
		Convey(fmt.Sprintf("Estimated document count should be %d", len(e_mocks.Users)), func() {
			count := LocalUserModel.EstimatedDocumentCount()
			So(count, ShouldEqual, len(e_mocks.Users))
		})
		Convey("Stats", func() {
			Convey("As a whole", func() {
				stats := LocalUserModel.Stats()
				So(stats.Count, ShouldEqual, len(e_mocks.Users))
				So(stats.AvgObjSize, ShouldBeGreaterThan, 0)
				So(stats.Size, ShouldBeGreaterThan, 0)
				So(stats.StorageSize, ShouldBeGreaterThan, 0)
			})
			Convey("Just the total size", func() {
				size := LocalUserModel.TotalSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the storage size", func() {
				size := LocalUserModel.StorageSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the total index size", func() {
				size := LocalUserModel.TotalIndexSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the average object size", func() {
				size := LocalUserModel.AvgObjSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the index count", func() {
				size := LocalUserModel.NumberOfIndexes()
				So(size, ShouldBeGreaterThanOrEqualTo, 1)
			})
		})
	})
}
