package elemental

import (
	"context"
	"github.com/elcengine/elemental/utils"
	"github.com/samber/lo"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CtxUser = "user"

type AuditType string

const (
	AuditTypeInsert AuditType = "insert"
	AuditTypeUpdate AuditType = "update"
	AuditTypeDelete AuditType = "delete"
)

type Audit struct {
	Entity    string      `json:"entity" bson:"entity"`         // The name of the model that was audited.
	Type      AuditType   `json:"type" bson:"type"`             // The type of operation that was performed (insert, update, delete).
	Document  primitive.M `json:"document" bson:"document"`     // The document that was affected by the operation.
	User      string      `json:"user" bson:"user"`             // The user who performed the operation if available within the context.
	CreatedAt time.Time   `json:"created_at" bson:"created_at"` // The date and time when the operation was performed.
}

var AuditModel = NewModel[Audit]("Audit", NewSchema(map[string]Field{
	"Entity": {
		Type:     reflect.String,
		Required: true,
	},
	"Type": {
		Type: reflect.String,
	},
	"Document": {
		Type:     reflect.Map,
		Required: true,
	},
	"User": {
		Type: reflect.String,
	},
}, SchemaOptions{
	Collection:              "audits",
	BypassSchemaEnforcement: true,
}))

// Enables auditing for the current model.
func (m Model[T]) EnableAuditing(ctx ...context.Context) {
	context := utils.CtxOrDefault(ctx)

	execWithModelOpts := func(q Model[Audit]) {
		if m.temporaryConnection != nil {
			q = q.SetConnection(*m.temporaryConnection)
		}
		if m.temporaryDatabase != nil {
			q = q.SetDatabase(*m.temporaryDatabase)
		}
		q.Exec(context)
	}

	userFallback := "System"

	m.OnInsert(func(doc T) {
		execWithModelOpts(AuditModel.Create(Audit{
			Entity:    m.Name,
			Type:      AuditTypeInsert,
			Document:  *utils.ToBSONDoc(doc),
			User:      lo.CoalesceOrEmpty(utils.Cast[string](context.Value(CtxUser)), userFallback),
			CreatedAt: time.Now(),
		}))
	}, TriggerOptions{Context: &context, Filter: &primitive.M{"ns.coll": primitive.M{"$eq": m.Collection().Name()}}})
	m.OnUpdate(func(doc T) {
		execWithModelOpts(AuditModel.Create(Audit{
			Entity:    m.Name,
			Type:      AuditTypeUpdate,
			Document:  *utils.ToBSONDoc(doc),
			User:      lo.CoalesceOrEmpty(utils.Cast[string](context.Value(CtxUser)), userFallback),
			CreatedAt: time.Now(),
		}))
	}, TriggerOptions{Context: &context, Filter: &primitive.M{"ns.coll": primitive.M{"$eq": m.Collection().Name()}}})
	m.OnDelete(func(id primitive.ObjectID) {
		execWithModelOpts(AuditModel.Create(Audit{
			Entity:    m.Name,
			Type:      AuditTypeDelete,
			Document:  map[string]any{"_id": id},
			User:      lo.CoalesceOrEmpty(utils.Cast[string](context.Value(CtxUser)), userFallback),
			CreatedAt: time.Now(),
		}))
	}, TriggerOptions{Context: &context, Filter: &primitive.M{"ns.coll": primitive.M{"$eq": m.Collection().Name()}}})
}
