package e_tests

import (
	"reflect"
	"testing"

	"github.com/akalanka47000/go-modkit/parallel_convey"
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/setup"
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

	AuditModel := elemental.AuditModel.SetDatabase(t.Name())

	ParallelConvey, Wait := pc.New(t)

	ParallelConvey("Insert", t, func() {
		KingdomModel.Create(Kingdom{Name: "Nilfgaard"}).Exec()
		SoTimeout(t, func() (ok bool) {
			audit := AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeInsert}).ExecT()
			if audit.Type != "" {
				ok = true
			}
			return
		})
	})

	ParallelConvey("Update", t, func() {
		KingdomModel.UpdateOne(&primitive.M{"name": "Nilfgaard"}, Kingdom{Name: "Redania"}).Exec()
		SoTimeout(t, func() (ok bool) {
			audit := AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeUpdate}).ExecT()
			if audit.Type != "" {
				ok = true
			}
			return
		})
	})

	ParallelConvey("Delete", t, func() {
		KingdomModel.Create(Kingdom{Name: "Skellige"}).Exec()
		KingdomModel.DeleteOne(primitive.M{"name": "Skellige"}).Exec()
		SoTimeout(t, func() (ok bool) {
			audit := AuditModel.FindOne(primitive.M{"entity": entity, "type": elemental.AuditTypeDelete}).ExecT()
			if audit.Type != "" {
				ok = true
			}
			return
		})
	})

	Wait()
}
