package e_tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/elcengine/elemental/utils"
	"github.com/akalanka47000/go-modkit/parallel_convey"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreAudit(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	entity := "Kingdom-For-Audit"

	KingdomModel := elemental.NewModel[Kingdom](entity, elemental.NewSchema(map[string]elemental.Field{
		"Name": {
			Type:     reflect.String,
			Required: true,
		},
	})).SetDatabase(t.Name())

	KingdomModel.EnableAuditing()

	Convey("Inspect audit records", t, func() {
		ParallelConvey, Wait := pc.New(t)

		ParallelConvey("Insert", func() {
			KingdomModel.Create(Kingdom{Name: "Nilfgaard"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](elemental.AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeInsert}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			}, time.After(15*time.Second))
		})

		ParallelConvey("Update", func() {
			KingdomModel.UpdateOne(&primitive.M{"name": "Nilfgaard"}, Kingdom{Name: "Redania"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](elemental.AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeUpdate}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			}, time.After(15*time.Second))
		})
		
		ParallelConvey("Delete", func() {
			KingdomModel.DeleteOne(primitive.M{"name": "Redania"}).Exec()
			SoTimeout(t, func() (ok bool) {
				audit := e_utils.Cast[elemental.Audit](elemental.AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeDelete}).Exec())
				if audit.Type != "" {
					ok = true
				}
				return
			}, time.After(15*time.Second))
		})

		Wait()
	})
}
