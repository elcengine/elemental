package e_tests

import (
	"os"
	"reflect"
	"testing"

	"github.com/akalanka47000/go-modkit/parallel_convey"
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/elcengine/elemental/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreAudit(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in non-CI environment")
	}

	e_test_setup.Connection(t.Name())

	entity := "Kingdom-For-Audit"

	KingdomModel := elemental.NewModel[Kingdom](entity, elemental.NewSchema(map[string]elemental.Field{
		"Name": {
			Type:     reflect.String,
			Required: true,
		},
	})).SetDatabase(t.Name())

	KingdomModel.EnableAuditing()

	AuditModel := elemental.AuditModel.SetDatabase(t.Name())

	Convey("Inspect audit records", t, func() {
		ParallelConvey, Wait := pc.New(t)

		ParallelConvey("Insert", func() {
			KingdomModel.Create(Kingdom{Name: "Nilfgaard"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeInsert}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			})
		})

		ParallelConvey("Update", func() {
			KingdomModel.UpdateOne(&primitive.M{"name": "Nilfgaard"}, Kingdom{Name: "Redania"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeUpdate}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			})
		})

		ParallelConvey("Delete", func() {
			KingdomModel.Create(Kingdom{Name: "Skellige"}).Exec()
			KingdomModel.DeleteOne(primitive.M{"name": "Skellige"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeDelete}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			})
		})

		Wait()
	})
}
