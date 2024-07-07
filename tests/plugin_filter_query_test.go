package e_tests

import (
	filter_query "elemental/plugins/filter-query"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPluginFilterQuery(t *testing.T) {

	Convey("Filters", t, func() {
		Convey("Basic Filters", func() {
			Convey("Equality", func() {
				So(1, ShouldEqual, 1)
			})
		})
	})

	Convey("Include", t, func() {
		Convey("When present in query string", func() {
			result := filter_query.Parse("include=field1,field2")
			So(result.Include, ShouldResemble, []string{"field1", "field2"})
		})
		Convey("When not present in query string", func() {
			result := filter_query.Parse("")
			So(len(result.Include), ShouldEqual, 0)
		})
	})

	Convey("Select", t, func() {
		Convey("When present in query string", func() {
			result := filter_query.Parse("select=field1,field2")
			So(result.Select, ShouldResemble, []string{"field1", "field2"})
		})
		Convey("When not present in query string", func() {
			result := filter_query.Parse("")
			So(len(result.Select), ShouldEqual, 0)
		})
	})

	Convey("Prepaginate", t, func() {
		Convey("When present in query string as true", func() {
			result := filter_query.Parse("prepaginate=true")
			So(result.Prepaginate, ShouldBeTrue)
		})
		Convey("When present in query string as false", func() {
			result := filter_query.Parse("prepaginate=false")
			So(result.Prepaginate, ShouldBeFalse)
		})
		Convey("When not present in query string", func() {
			result := filter_query.Parse("")
			So(result.Prepaginate, ShouldBeFalse)
		})
	})
}
