package e_tests

import (
	"testing"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {

	e_test_setup.Connection()

	defer e_test_setup.Teardown()

	Convey("Batch transaction", t, func() {
		Convey("Between 2 databases within the same connection", func() {
			Convey("Should be able to insert into both databases", func() {
				elemental.TransactionBatch(
					UserModel.Create(User{
						Name: "Yennefer",
					}).FlexibleClone(),
					UserModel.Create(User{
						Name: "Triss",
					}).SetDatabase(e_mocks.SECONDARY_DB).FlexibleClone(),
				)
				yennefer := UserModel.FindOne().Where("name", "Yennefer").Exec()
				So(yennefer, ShouldNotBeNil)
				triss := UserModel.FindOne().SetDatabase(e_mocks.SECONDARY_DB).Where("name", "Triss").Exec()
				So(triss, ShouldNotBeNil)
			})
			// Convey("Should rollback if one of the operations fail", func() {
			// 	elemental.TransactionBatch(
			// 		UserModel.Create(User{
			// 			Name: "Eskel",
			// 		}).FlexibleClone(),
			// 		UserModel.Create(User{
			// 			Name: "Eredin",
			// 		}).SetDatabase(e_mocks.SECONDARY_DB).FlexibleClone(),
			// 		UserModel.Create(User{
			// 			Name: "Eredin",
			// 		}).SetDatabase(e_mocks.SECONDARY_DB).FlexibleClone(),
			// 	)
			// 	eskel := UserModel.FindOne().Where("name", "Eskel").Exec().(User)
			// 	So(eskel, ShouldBeNil)
			// 	eredin := UserModel.FindOne().SetDatabase(e_mocks.SECONDARY_DB).Where("name", "Eredin").Exec().(User)
			// 	So(eredin, ShouldBeNil)
			// })
		})
	})
}
