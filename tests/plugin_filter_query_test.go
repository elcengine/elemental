//nolint:maintidx
package tests

import (
	"fmt"
	"testing"
	"time"

	fq "github.com/elcengine/elemental/plugins/filterquery"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/samber/lo"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPluginFilterQuery(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Filters", t, func() {
		Convey("Basic Syntax", func() {
			Convey("Equality", func() {
				result := fq.Parse("filter[name]=John")
				So(result.Filters, ShouldResemble, bson.M{"name": "John"})
			})
			Convey("Equality with a boolean", func() {
				result := fq.Parse("filter[active]=true")
				So(result.Filters, ShouldResemble, bson.M{"active": true})
			})
		})
		Convey("Advanced Syntax", func() {
			Convey("Equality", func() {
				result := fq.Parse("filter[name]=eq(John)")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$eq": "John"}})
			})
			Convey("Inequality", func() {
				result := fq.Parse("filter[age]=ne(30)")
				So(result.Filters, ShouldResemble, bson.M{"age": bson.M{"$ne": 30.0}})
			})
			Convey("Equality with a boolean", func() {
				result := fq.Parse("filter[active]=eq(true)")
				So(result.Filters, ShouldResemble, bson.M{"active": bson.M{"$eq": true}})
			})
			Convey("Equality with ObjectId", func() {
				id := primitive.NewObjectID()
				result := fq.Parse(fmt.Sprintf("filter[_id]=eq(%s)", id.Hex()))
				So(result.Filters, ShouldResemble, bson.M{"_id": bson.M{"$eq": id}})
			})
			Convey("Greater Than", func() {
				result := fq.Parse("filter[age]=gt(30)")
				So(result.Filters, ShouldResemble, bson.M{"age": bson.M{"$gt": 30.0}})
			})
			Convey("Greater Than with Time", func() {
				result := fq.Parse("filter[createdAt]=gt(2023-01-01T00:00:00Z)")
				So(result.Filters, ShouldResemble, bson.M{"createdAt": bson.M{"$gt": lo.Must(time.Parse(time.RFC3339, "2023-01-01T00:00:00Z"))}})
			})
			Convey("Less Than", func() {
				result := fq.Parse("filter[age]=lt(30)")
				So(result.Filters, ShouldResemble, bson.M{"age": bson.M{"$lt": 30.0}})
			})
			Convey("Greater Than or Equal", func() {
				result := fq.Parse("filter[age]=gte(30)")
				So(result.Filters, ShouldResemble, bson.M{"age": bson.M{"$gte": 30.0}})
			})
			Convey("Less Than or Equal", func() {
				result := fq.Parse("filter[age]=lte(30)")
				So(result.Filters, ShouldResemble, bson.M{"age": bson.M{"$lte": 30.0}})
			})
			Convey("In", func() {
				result := fq.Parse("filter[name]=in(John,Jane)")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$in": []any{"John", "Jane"}}})
			})
			Convey("Not In", func() {
				result := fq.Parse("filter[name]=nin(John,Jane)")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$nin": []any{"John", "Jane"}}})
			})
			Convey("Regex", func() {
				result := fq.Parse("filter[name]=reg(^J)")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^J", Options: ""}}})
			})
			Convey("Regex with modifiers", func() {
				result := fq.Parse("filter[name]=reg(^J...[i])")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^J", Options: "i"}}})
			})
			Convey("Exists", func() {
				result := fq.Parse("filter[name]=exists(true)")
				So(result.Filters, ShouldResemble, bson.M{"name": bson.M{"$exists": true}})
			})
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(len(result.Filters), ShouldEqual, 0)
		})
		Convey("Complex operators", func() {
			Convey("Or", func() {
				result := fq.Parse("filter[or]=name=eq(John),age=gt(30)")
				So(result.Filters, ShouldResemble, bson.M{"$or": []any{
					bson.M{"name": bson.M{"$eq": "John"}},
					bson.M{"age": bson.M{"$gt": 30.0}},
				}})
			})
			Convey("Or within values", func() {
				result := fq.Parse("filter[age]=or(lt(30),gt(40))")
				So(result.Filters, ShouldResemble, bson.M{"$or": []any{
					bson.M{"age": bson.M{"$lt": 30.0}},
					bson.M{"age": bson.M{"$gt": 40.0}},
				}})
			})
			Convey("And", func() {
				result := fq.Parse("filter[and]=name=eq(John),age=gt(30)")
				So(result.Filters, ShouldResemble, bson.M{"$and": []any{
					bson.M{"name": bson.M{"$eq": "John"}},
					bson.M{"age": bson.M{"$gt": 30.0}},
				}})
			})
		})
	})

	Convey("Sorts", t, func() {
		Convey("Ascending", func() {
			result := fq.Parse("sort[name]=asc")
			So(result.Sorts, ShouldResemble, bson.M{"name": 1})
		})
		Convey("Ascending with 1", func() {
			result := fq.Parse("sort[name]=1")
			So(result.Sorts, ShouldResemble, bson.M{"name": 1})
		})
		Convey("Descending", func() {
			result := fq.Parse("sort[name]=desc")
			So(result.Sorts, ShouldResemble, bson.M{"name": -1})
		})
		Convey("Descending with -1", func() {
			result := fq.Parse("sort[name]=-1")
			So(result.Sorts, ShouldResemble, bson.M{"name": -1})
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(len(result.Sorts), ShouldEqual, 0)
		})
	})

	Convey("Include", t, func() {
		Convey("When present in query string", func() {
			result := fq.Parse("include=field1,field2")
			So(result.Include, ShouldResemble, []string{"field1", "field2"})
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(len(result.Include), ShouldEqual, 0)
		})
	})

	Convey("Select", t, func() {
		Convey("When present in query string", func() {
			result := fq.Parse("select=field1,field2")
			So(result.Select, ShouldResemble, bson.M{"field1": 1, "field2": 1})
		})
		Convey("When present in query string with exclusion", func() {
			result := fq.Parse("select=-field1,field2")
			So(result.Select, ShouldResemble, bson.M{"field1": 0, "field2": 1})
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(len(result.Select), ShouldEqual, 0)
		})
	})

	Convey("Prepaginate", t, func() {
		Convey("When present in query string as true", func() {
			result := fq.Parse("prepaginate=true")
			So(result.Prepaginate, ShouldBeTrue)
		})
		Convey("When present in query string as false", func() {
			result := fq.Parse("prepaginate=false")
			So(result.Prepaginate, ShouldBeFalse)
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(result.Prepaginate, ShouldBeFalse)
		})
	})

	Convey("Page", t, func() {
		Convey("When present in query string", func() {
			result := fq.Parse("page=2")
			So(result.Page, ShouldEqual, 2)
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(result.Page, ShouldEqual, 0)
		})
		Convey("When present in query string with invalid value", func() {
			result := fq.Parse("page=invalid")
			So(result.Page, ShouldEqual, 0)
		})
		Convey("When present in query string with negative value", func() {
			result := fq.Parse("page=-1")
			So(result.Page, ShouldEqual, 0)
		})
	})

	Convey("Limit", t, func() {
		Convey("When present in query string", func() {
			result := fq.Parse("limit=10")
			So(result.Limit, ShouldEqual, 10)
		})
		Convey("When not present in query string", func() {
			result := fq.Parse("")
			So(result.Limit, ShouldEqual, 0)
		})
		Convey("When present in query string with invalid value", func() {
			result := fq.Parse("limit=invalid")
			So(result.Limit, ShouldEqual, 0)
		})
		Convey("When present in query string with negative value", func() {
			result := fq.Parse("limit=-10")
			So(result.Limit, ShouldEqual, 0)
		})
	})

	Convey("QS", t, func() {
		Convey("When a filter is present in query string", func() {
			results := UserModel.QS(fmt.Sprintf("filter[name]=eq(%s)", mocks.Caranthir.Name)).ExecTT()
			So(results, ShouldHaveLength, 1)
			So(results[0].Name, ShouldEqual, mocks.Caranthir.Name)
		})
		Convey("When a secondary filter is present in query string", func() {
			results := UserModel.QS(fmt.Sprintf("secondaryFilter[name]=eq(%s)", mocks.Vesemir.Name)).ExecTT()
			So(results, ShouldHaveLength, 1)
			So(results[0].Name, ShouldEqual, mocks.Vesemir.Name)
		})
		Convey("When a sort is present in query string", func() {
			results := UserModel.QS("sort[name]=desc").ExecTT()
			So(results, ShouldHaveLength, len(mocks.Users))
			So(results[0].Name, ShouldEqual, mocks.Yennefer.Name)
		})
		Convey("When a select is present in query string", func() {
			results := UserModel.QS(fmt.Sprintf("select=age&filter[name]=eq(%s)", mocks.Geralt.Name)).ExecTT()
			So(results, ShouldHaveLength, 1)
			So(results[0].ID, ShouldNotBeZeroValue)
			So(results[0].Name, ShouldBeZeroValue)
			So(results[0].Age, ShouldEqual, mocks.Geralt.Age)
		})
		Convey("When a page and limit are present in query string", func() {
			result := UserModel.QS("page=1&limit=2").ExecTP()
			So(result.Docs, ShouldHaveLength, 2)
			So(result.Docs[0].Name, ShouldEqual, mocks.Ciri.Name)
			So(result.Docs[1].Name, ShouldEqual, mocks.Geralt.Name)
		})
		Convey("When an include is present in query string", func() {
			MonsterModel := MonsterModel.SetDatabase(t.Name())
			KingdomModel := KingdomModel.SetDatabase(t.Name())
			BestiaryModel := BestiaryModel.SetDatabase(t.Name())

			monster := MonsterModel.Create(Monster{
				Name:     "Katakan",
				Category: "Vampire",
			}).ExecT()

			kingdom := KingdomModel.Create(Kingdom{
				Name: "Nilfgaard",
			}).ExecT()

			BestiaryModel.Create(Bestiary{
				Monster: monster,
				Kingdom: kingdom,
			}).Exec()

			bestiaries := BestiaryModel.QS("include=monster,kingdom").ExecTT()
			So(bestiaries, ShouldHaveLength, 1)
			So(bestiaries[0].Monster.Name, ShouldEqual, monster.Name)
			So(bestiaries[0].Monster.Category, ShouldEqual, monster.Category)
			So(bestiaries[0].Kingdom.Name, ShouldEqual, kingdom.Name)
		})
	})
}
