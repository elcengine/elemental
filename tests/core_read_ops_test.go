package e_tests

import (
	"fmt"
	"testing"

	e_test_base "github.com/elcengine/elemental/tests/base"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreReadOps(t *testing.T) {
	t.Parallel()

	e_test_setup.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Read users with operators", t, func() {
		Convey(fmt.Sprintf("Find all where age is %d", e_mocks.Geralt.Age), func() {
			users := UserModel.Where("age").Equals(e_mocks.Geralt.Age).ExecTT()
			So(len(users), ShouldEqual, 2)
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find all where age is greater than 50", func() {
			users := UserModel.Where("age").GreaterThan(50).ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Age > 50
			})))
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			So(users[2].Name, ShouldEqual, e_mocks.Imlerith.Name)
			So(users[3].Name, ShouldEqual, e_mocks.Yennefer.Name)
			So(users[4].Name, ShouldEqual, e_mocks.Vesemir.Name)
		})
		Convey("Find a mage where age is greater than 50", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"occupation": "Mage"}).Where("age").GreaterThan(50).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
			Convey("In conjuntion with find one", func() {
				user := UserModel.FindOne(primitive.M{"occupation": "Mage"}).Where("age").GreaterThan(50).ExecPtr()
				So(user, ShouldNotBeNil)
				So(user.Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
			Convey("In conjuntion with equals", func() {
				users := UserModel.Where("age").GreaterThan(50).Where("occupation").Equals("Mage").ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
		})
		Convey("Find where age is between 90 and 110", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"$and": []primitive.M{
					{"age": primitive.M{"$gte": 90}},
					{"age": primitive.M{"$lte": 110}},
				}}).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
			Convey("In conjuntion with where", func() {
				users := UserModel.Where("age").GreaterThanOrEquals(90).Where("age").LessThanOrEquals(110).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
			Convey("In conjuntion with between", func() {
				users := UserModel.Where("age").Between(90, 110).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
		})
		Convey("Find where age is 120 or 150", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"$or": []primitive.M{
					{"age": 120},
					{"age": 150},
				}}).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
			Convey("In conjuntion with in", func() {
				users := UserModel.Where("age").In(120, 150).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
			Convey("In conjuntion with or operator", func() {
				users := UserModel.Where("age").Equals(120).Or().Where("age").Equals(150).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
			Convey("In conjuntion with or where operator", func() {
				users := UserModel.Where("age").Equals(120).OrWhere("age").Equals(150).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
		})
		Convey(fmt.Sprintf("Find where age is not %d", e_test_base.DefaultAge), func() {
			expectedCount := len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Age > 0 && u.Age != e_test_base.DefaultAge
			}))
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"age": primitive.M{"$ne": e_test_base.DefaultAge}}).ExecTT()
				So(len(users), ShouldEqual, expectedCount)
			})
			Convey("In conjuntion with not equals", func() {
				users := UserModel.Where("age").NotEquals(e_test_base.DefaultAge).ExecTT()
				So(len(users), ShouldEqual, expectedCount)
			})
			Convey("In conjuntion with not in", func() {
				users := UserModel.Where("age").NotIn(e_test_base.DefaultAge).ExecTT()
				So(len(users), ShouldEqual, expectedCount)
			})
		})
		Convey("Find where weapon list contains Battle Axe", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"weapons": "Battle Axe"}).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
			Convey("In conjuntion with element match", func() {
				users := UserModel.Where("weapons").ElementMatches(primitive.M{"$eq": "Battle Axe"}).ExecTT()
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			})
		})
		Convey(fmt.Sprintf("Find where weapon count is %d", len(e_mocks.Geralt.Weapons)), func() {
			users := UserModel.Where("weapons").Size(len(e_mocks.Geralt.Weapons)).ExecTT()
			So(len(users), ShouldEqual, 1)
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find where occupation exists", func() {
			users := UserModel.Where("occupation").Exists(true).ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Occupation != ""
			})))
		})
		Convey("Find where occupation does not exist", func() {
			users := UserModel.Where("occupation").Exists(false).ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Occupation == ""
			})))
		})
		Convey("Find where name matches the pattern", func() {
			users := UserModel.Where("name").Regex(".*alt").ExecTT()
			So(len(users), ShouldEqual, 1)
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find where age is divisible by 5", func() {
			users := UserModel.Where("age").Mod(5, 0).ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Age > 0 && u.Age%5 == 0
			})))
		})
		Convey("Count users in conjuntion with greater than", func() {
			count := UserModel.Where("age").GreaterThan(50).CountDocuments().ExecInt()
			So(count, ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.Age > 50
			})))
		})
	})
}
