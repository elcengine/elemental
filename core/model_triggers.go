package elemental

import (
	"context"

	"github.com/elcengine/elemental/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type triggerType string

const (
	triggerTypeInsert           triggerType = "insert"
	triggerTypeUpdate           triggerType = "update"
	triggerTypeDelete           triggerType = "delete"
	triggerTypeReplace          triggerType = "replace"
	triggerTypeCollectionDrop   triggerType = "drop"
	triggerTypeCollectionRename triggerType = "rename"
	triggerTypeStreamInvalidate triggerType = "invalidate"
)

type TriggerOptions struct {
	Filter  *primitive.M     // Optional filter to apply to the change stream
	Context *context.Context // Optional context to use for the change stream. If not provided, the default context will be used.
}

func (m Model[T]) on(event triggerType, f func(change primitive.M), opts ...TriggerOptions) *mongo.ChangeStream {
	filters := primitive.M{}
	ctx := context.Background()
	if len(opts) > 0 {
		if opts[0].Context != nil {
			ctx = *opts[0].Context
		}
		if opts[0].Filter != nil {
			filters = *opts[0].Filter
		}
	}
	filters["operationType"] = event
	cs, err := m.Collection().Watch(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: filters}},
	}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		panic(err)
	}
	go func() {
		for cs.Next(ctx) {
			var changeDoc bson.M
			if err := cs.Decode(&changeDoc); err != nil {
				cs.Close(ctx)
				return
			}
			f(changeDoc)
		}
	}()
	go func() {
		<-m.triggerExit
		cs.Close(ctx)
	}()
	return cs
}

// Listens for insert events on the collection and calls the provided function with the inserted document.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/insert
func (m Model[T]) OnInsert(f func(doc T), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeInsert, func(change primitive.M) {
		f(utils.CastBSON[T](change["fullDocument"]))
	}, opts...)
}

// Listens for update events on the collection and calls the provided function with the updated document.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/update
func (m Model[T]) OnUpdate(f func(doc T), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeUpdate, func(change primitive.M) {
		f(utils.CastBSON[T](change["fullDocument"]))
	}, opts...)
}

// Listens for delete events on the collection and calls the provided function with the deleted document's ID.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/delete
func (m Model[T]) OnDelete(f func(id primitive.ObjectID), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeDelete, func(change primitive.M) {
		f(utils.Cast[primitive.ObjectID](change["documentKey"].(primitive.M)["_id"]))
	}, opts...)
}

// Listens for replace events on the collection and calls the provided function with the replaced document.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/replace
func (m Model[T]) OnReplace(f func(doc T), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeReplace, func(change primitive.M) {
		f(utils.CastBSON[T](change["fullDocument"]))
	}, opts...)
}

// Listens for collection drop events and calls the provided function.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/drop
func (m Model[T]) OnCollectionDrop(f func(), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeCollectionDrop, func(change primitive.M) {
		f()
	}, opts...)
}

// Listens for collection rename events and calls the provided function.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/rename
func (m Model[T]) OnCollectionRename(f func(), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeCollectionRename, func(change primitive.M) {
		f()
	}, opts...)
}

// Listens for stream invalidation events and calls the provided function.
// For more information, refer the following link: https://www.mongodb.com/docs/manual/reference/change-events/invalidate
func (m Model[T]) OnStreamInvalidate(f func(), opts ...TriggerOptions) *mongo.ChangeStream {
	return m.on(triggerTypeStreamInvalidate, func(change primitive.M) {
		f()
	}, opts...)
}

// InvalidateTriggers clears all triggers for the model. This includes audits
func (m Model[T]) InvalidateTriggers() {
	m.triggerExit <- true
}
