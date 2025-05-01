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

	SECONDARY_DB := fmt.Sprintf("%s_%s", t.Name(), "secondary")

	UserModel.SyncIndexes()

	Convey("Batch transaction", t, func() {
		Convey("Between 2 databases within the same connection", func() {
			Convey("Should be able to insert into both databases", func() {
				elemental.TransactionBatch(
					UserModel.Create(User{
						Name: "Yennefer",
					}),
					UserModel.SetDatabase(SECONDARY_DB).Create(User{
						Name: "Triss",
					}),
				)
				yennefer := UserModel.FindOne().Where("name", "Yennefer").Exec()
				So(yennefer, ShouldNotBeNil)
				triss := UserModel.SetDatabase(SECONDARY_DB).FindOne().Where("name", "Triss").Exec()
				So(triss, ShouldNotBeNil)
			})
			Convey("Should rollback if one of the operations fail", func() {
				elemental.TransactionBatch(
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
				eskel := UserModel.FindOne().Where("name", "Eskel").Exec()
				So(eskel, ShouldBeNil)
				eredin := UserModel.FindOne().SetDatabase(SECONDARY_DB).Where("name", "Eredin").Exec()
				So(eredin, ShouldBeNil)
			})
		})
	})
}
