package tests

import (
	"fmt"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

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

	Convey("Basic transaction", t, func() {
		result, err := elemental.Transaction(func(ctx mongo.SessionContext) (any, error) {
			return UserModel.Create(User{
				Name: uuid.NewString(),
			}).Exec(ctx), nil
		})
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})

	Convey("Client transaction", t, func() {
		alias := uuid.NewString()
		elemental.Connect(elemental.ConnectionOptions{
			URI:   mocks.DEFAULT_DATASOURCE,
			Alias: alias,
		})
		result, err := elemental.ClientTransaction(alias, func(ctx mongo.SessionContext) (any, error) {
			return UserModel.Create(User{
				Name: uuid.NewString(),
			}).SetConnection(alias).Exec(ctx), nil
		})
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})
}
