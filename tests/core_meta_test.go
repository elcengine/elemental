package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreMeta(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	UserModel.SyncIndexes()

	UserModel.InsertMany(mocks.Users).Exec()

	Convey("Retrieve the underlying client used by a model", t, func() {
		client := UserModel.Client()
		So(client, ShouldNotBeNil)
		So(client.Ping(context.Background(), nil), ShouldBeNil)
	})

	Convey("Metadata", t, func() {
		Convey(fmt.Sprintf("Estimated document count should be %d", len(mocks.Users)), func() {
			count := UserModel.EstimatedDocumentCount()
			So(count, ShouldEqual, len(mocks.Users))
		})
		Convey("Stats", func() {
			Convey("As a whole", func() {
				stats := UserModel.Stats()
				So(stats.Count, ShouldEqual, len(mocks.Users))
				So(stats.AvgObjSize, ShouldBeGreaterThan, 0)
				So(stats.Size, ShouldBeGreaterThan, 0)
				So(stats.StorageSize, ShouldBeGreaterThan, 0)
			})
			Convey("Just the total size", func() {
				size := UserModel.TotalSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the storage size", func() {
				size := UserModel.StorageSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the total index size", func() {
				size := UserModel.TotalIndexSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the average object size", func() {
				size := UserModel.AvgObjSize()
				So(size, ShouldBeGreaterThan, 0)
			})
			Convey("Just the index count", func() {
				size := UserModel.NumberOfIndexes()
				So(size, ShouldEqual, 2)
			})
		})
	})
}
