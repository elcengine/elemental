package tests

import (
	"fmt"
	"testing"

	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreReadSelect(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Find with only specified fields", t, func() {
		limit := int64(2)
		Convey(fmt.Sprintf("%d user names with ID", limit), func() {
			assert := func(users []User) {
				So(users, ShouldHaveLength, limit)
				for _, user := range users {
					So(user.ID, ShouldNotBeZeroValue)
					So(user.Name, ShouldNotBeZeroValue)
					So(user.School, ShouldBeNil)
					So(user.CreatedAt, ShouldBeZeroValue)
				}
			}
			Convey("In conjunction with a map input", func() {
				users := UserModel.Find().Select(primitive.M{"name": 1}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a slice input", func() {
				users := UserModel.Find().Select([]string{"name"}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a string input", func() {
				users := UserModel.Find().Select("name").Limit(limit).ExecTT()
				assert(users)
			})
		})
		Convey(fmt.Sprintf("%d user names without ID", limit), func() {
			assert := func(users []User) {
				So(users, ShouldHaveLength, limit)
				for _, user := range users {
					So(user.ID, ShouldBeZeroValue)
					So(user.Name, ShouldNotBeZeroValue)
					So(user.School, ShouldBeNil)
					So(user.CreatedAt, ShouldBeZeroValue)
				}
			}
			Convey("In conjunction with a map input", func() {
				users := UserModel.Find().Select(primitive.M{"name": 1, "_id": 0}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a slice input", func() {
				users := UserModel.Find().Select([]string{"name", "-_id"}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a string input (spaces)", func() {
				users := UserModel.Find().Select("name -_id").Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a string input (commas)", func() {
				users := UserModel.Find().Select("name, -_id").Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with variadic arguments", func() {
				users := UserModel.Find().Select("name", "-_id").Limit(limit).ExecTT()
				assert(users)
			})
		})
		Convey(fmt.Sprintf("%d user names and ages", limit), func() {
			assert := func(users []User) {
				So(users, ShouldHaveLength, limit)
				for _, user := range users {
					So(user.ID, ShouldNotBeZeroValue)
					So(user.Name, ShouldNotBeZeroValue)
					So(user.Age, ShouldNotBeZeroValue)
				}
			}
			Convey("In conjunction with a map input", func() {
				users := UserModel.Find().Select(primitive.M{"name": 1, "age": 1}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a slice input", func() {
				users := UserModel.Find().Select([]string{"name", "age"}).Limit(limit).ExecTT()
				assert(users)
			})
			Convey("In conjunction with a string input", func() {
				users := UserModel.Find().Select("name age").Limit(limit).ExecTT()
				assert(users)
			})
		})
	})
}
