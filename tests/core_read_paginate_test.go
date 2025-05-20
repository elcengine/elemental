package tests

import (
	"testing"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreReadPaginate(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Find paginated users", t, func() {
		Convey("First page", func() {
			result := UserModel.Find().Paginate(1, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 2)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 1)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(mocks.Users))
			So(result.HasPrev, ShouldBeFalse)
			So(result.HasNext, ShouldBeTrue)
			So(result.NextPage, ShouldEqual, lo.ToPtr[int64](2))
			So(result.PrevPage, ShouldBeNil)
			So(result.Docs[0].Name, ShouldEqual, mocks.Ciri.Name)
			So(result.Docs[1].Name, ShouldEqual, mocks.Geralt.Name)
		})
		Convey("Second page", func() {
			result := UserModel.Find().Paginate(2, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 2)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 2)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(mocks.Users))
			So(result.HasPrev, ShouldBeTrue)
			So(result.HasNext, ShouldBeTrue)
			So(result.NextPage, ShouldEqual, lo.ToPtr[int64](3))
			So(result.PrevPage, ShouldEqual, lo.ToPtr[int64](1))
			So(result.Docs[0].Name, ShouldEqual, mocks.Eredin.Name)
			So(result.Docs[1].Name, ShouldEqual, mocks.Caranthir.Name)
		})
		Convey("Last page", func() {
			result := UserModel.Find().Paginate(4, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 1)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 4)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(mocks.Users))
			So(result.HasPrev, ShouldBeTrue)
			So(result.HasNext, ShouldBeFalse)
			So(result.NextPage, ShouldBeNil)
			So(result.PrevPage, ShouldEqual, lo.ToPtr[int64](3))
			So(result.Docs[0].Name, ShouldEqual, mocks.Vesemir.Name)
		})
		Convey("First page with filters", func() {
			result := UserModel.Find(primitive.M{"name": "Ciri"}).Paginate(1, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 1)
			So(result.TotalPages, ShouldEqual, 1)
			So(result.Page, ShouldEqual, 1)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, 1)
			So(result.HasPrev, ShouldBeFalse)
			So(result.HasNext, ShouldBeFalse)
			So(result.NextPage, ShouldBeNil)
			So(result.PrevPage, ShouldBeNil)
			So(result.Docs[0].Name, ShouldEqual, mocks.Ciri.Name)
		})
		Convey("First page with filters (No results)", func() {
			result := UserModel.Find(primitive.M{"name": uuid.NewString()}).Paginate(1, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 0)
			So(result.TotalPages, ShouldEqual, 0)
			So(result.Page, ShouldEqual, 1)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, 0)
			So(result.HasPrev, ShouldBeFalse)
			So(result.HasNext, ShouldBeFalse)
			So(result.NextPage, ShouldBeNil)
			So(result.PrevPage, ShouldBeNil)
		})
	})
}
