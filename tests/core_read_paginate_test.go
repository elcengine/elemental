package e_tests

import (
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"testing"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreReadPaginate(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Find paginated users", t, func() {
		Convey("First page", func() {
			t.Parallel()
			result := UserModel.Find().Paginate(1, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 2)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 1)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(e_mocks.Users))
			So(result.HasPrev, ShouldBeFalse)
			So(result.HasNext, ShouldBeTrue)
			So(result.NextPage, ShouldEqual, lo.ToPtr[int64](2))
			So(result.PrevPage, ShouldBeNil)
			So(result.Docs[0].Name, ShouldEqual, e_mocks.Ciri.Name)
			So(result.Docs[1].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Second page", func() {
			t.Parallel()
			result := UserModel.Find().Paginate(2, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 2)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 2)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(e_mocks.Users))
			So(result.HasPrev, ShouldBeTrue)
			So(result.HasNext, ShouldBeTrue)
			So(result.NextPage, ShouldEqual, lo.ToPtr[int64](3))
			So(result.PrevPage, ShouldEqual, lo.ToPtr[int64](1))
			So(result.Docs[0].Name, ShouldEqual, e_mocks.Eredin.Name)
			So(result.Docs[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
		})
		Convey("Last page", func() {
			t.Parallel()
			result := UserModel.Find().Paginate(4, 2).Exec().(elemental.PaginateResult[User])
			So(len(result.Docs), ShouldEqual, 1)
			So(result.TotalPages, ShouldEqual, 4)
			So(result.Page, ShouldEqual, 4)
			So(result.Limit, ShouldEqual, 2)
			So(result.TotalDocs, ShouldEqual, len(e_mocks.Users))
			So(result.HasPrev, ShouldBeTrue)
			So(result.HasNext, ShouldBeFalse)
			So(result.NextPage, ShouldBeNil)
			So(result.PrevPage, ShouldEqual, lo.ToPtr[int64](3))
			So(result.Docs[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
		})
	})
}
