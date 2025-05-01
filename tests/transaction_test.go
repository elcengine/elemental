package e_tests

import (
	"testing"

	elemental "github.com/elcengine/elemental/core"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection()

	defer e_test_setup.Teardown()

	var TxnUserModel = UserModel.Clone().SetCollection("txn_users")

	TxnUserModel.SyncIndexes()

	Convey("Batch transaction", t, func() {
		Convey("Between 2 databases within the same connection", func() {
			Convey("Should be able to insert into both databases", func() {
				elemental.TransactionBatch(
					TxnUserModel.Create(User{
						Name: "Yennefer",
					}),
					TxnUserModel.SetDatabase(e_mocks.SECONDARY_DB).Create(User{
						Name: "Triss",
					}),
				)
				yennefer := TxnUserModel.FindOne().Where("name", "Yennefer").Exec()
				So(yennefer, ShouldNotBeNil)
				triss := TxnUserModel.SetDatabase(e_mocks.SECONDARY_DB).FindOne().Where("name", "Triss").Exec()
				So(triss, ShouldNotBeNil)
			})
			Convey("Should rollback if one of the operations fail", func() {
				elemental.TransactionBatch(
					TxnUserModel.Create(User{
						Name: "Eskel",
					}),
					TxnUserModel.Create(User{
						Name: "Eskel",
					}),
					TxnUserModel.SetDatabase(e_mocks.SECONDARY_DB).Create(User{
						Name: "Eredin",
					}),
				)
				eskel := TxnUserModel.FindOne().Where("name", "Eskel").Exec()
				So(eskel, ShouldBeNil)
				eredin := TxnUserModel.FindOne().SetDatabase(e_mocks.SECONDARY_DB).Where("name", "Eredin").Exec()
				So(eredin, ShouldBeNil)
			})
		})
	})
}
