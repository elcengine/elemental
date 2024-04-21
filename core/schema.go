package elemental

import (
	"context"
	"elemental/connection"
	"elemental/utils"
	"reflect"
	"github.com/creasty/defaults"
	"github.com/rxwycdh/rxhash"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Schema struct {
	Definitions map[string]Field
	Options     SchemaOptions
}

var schemas = make(map[string]Schema)

func NewSchema(definitions map[string]Field, opts SchemaOptions) Schema {
	key := lo.Must((rxhash.HashStruct(definitions)))
	if _, ok := schemas[key]; ok {
		return schemas[key]
	}
	defaults.Set(opts)
	schema := Schema{
		Definitions: definitions,
		Options:     opts,
	}
	schemas[key] = schema
	return schema
}

func (s Schema) syncIndexes(reflectedBaseType reflect.Type) {
	collection := e_connection.Use(s.Options.Database, s.Options.Connection).Collection(s.Options.Collection)
	collection.Indexes().DropAll(context.Background())
	for field, definition := range s.Definitions {
		if (definition.Index != options.IndexOptions{}) {
			reflectedField, _ := reflectedBaseType.FieldByName(field)
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: reflectedField.Tag.Get("bson") , Value: e_utils.Coalesce(definition.IndexOrder, 1)}},
				Options: &definition.Index,
			}
			collection.Indexes().CreateOne(context.TODO(), indexModel)
		}
	}
}
