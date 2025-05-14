package e_tests

import (
	"fmt"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	e_test_setup "github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	UserModel.SyncIndexes()

	SECONDARY_DB := fmt.Sprintf("%s_%s", t.Name(), "secondary")

	Convey("Batch transaction", t, func() {
		Convey("Between 2 databases within the same connection", func() {
			Convey("Should be able to insert into both databases", func() {
				results, errors := elemental.TransactionBatch(
					UserModel.Create(User{
						Name: "Yennefer",
					}),
					UserModel.Create(User{
						Name: "Triss",
					}).SetDatabase(SECONDARY_DB),
				)
				So(results, ShouldHaveLength, 2)
				So(errors, ShouldBeEmpty)
				yennefer := UserModel.FindOne().Where("name", "Yennefer").ExecPtr()
				So(yennefer, ShouldNotBeNil)
				triss := UserModel.FindOne().Where("name", "Triss").SetDatabase(SECONDARY_DB).ExecPtr()
				So(triss, ShouldNotBeNil)
			})
			Convey("Should rollback if one of the operations fail", func() {
				_, errors := elemental.TransactionBatch(
					UserModel.Create(User{
						Name: "Eskel",
					}),
					UserModel.Create(User{
						Name: "Eskel",
					}),
					UserModel.SetDatabase(SECONDARY_DB).Create(User{
						Name: "Eredin",
					}),
				)
				So(errors, ShouldNotBeEmpty)
				eskel := UserModel.FindOne().Where("name", "Eskel").ExecPtr()
				So(eskel, ShouldBeNil)
				eredin := UserModel.FindOne().SetDatabase(SECONDARY_DB).Where("name", "Eredin").ExecPtr()
				So(eredin, ShouldBeNil)
			})
		})
	})
}
