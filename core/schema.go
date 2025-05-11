package elemental

import (
	"context"
	"reflect"

	"github.com/creasty/defaults"
	"github.com/elcengine/elemental/connection"
	"github.com/elcengine/elemental/utils"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Schema struct {
	Definitions map[string]Field // The definitions of the schema, basically the fields of the document
	Options     SchemaOptions    // Custom schema options, like the collection name, database name, etc.
}

// Creates a new Elemental schema with the given definitions and options.
func NewSchema(definitions map[string]Field, opts ...SchemaOptions) Schema {
	schema := Schema{
		Definitions: definitions,
	}
	if len(opts) > 0 {
		defaults.Set(opts[0])
		schema.Options = opts[0]
	}
	return schema
}

// Retrives the Elemental field definition for the given path in the schema.
func (s Schema) Field(path string) *Field {
	definition := s.Definitions[path]
	if definition != (Field{}) {
		return &definition
	}
	return nil
}

func (s Schema) syncIndexes(reflectedBaseType reflect.Type, databaseOverride, connectionOverride, collectionOverride string) {
	database, _ := lo.Coalesce(databaseOverride, s.Options.Database)
	connection, _ := lo.Coalesce(connectionOverride, s.Options.Connection)
	collectionName, _ := lo.Coalesce(collectionOverride, s.Options.Collection)
	collection := e_connection.Use(database, connection).Collection(collectionName)
	collection.Indexes().DropAll(context.Background())
	for field, definition := range s.Definitions {
		if (definition.Index != options.IndexOptions{}) {
			reflectedField, _ := reflectedBaseType.FieldByName(field)
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: cleanBSONTag(reflectedField.Tag.Get("bson")), Value: e_utils.Coalesce(definition.IndexOrder, 1)}},
				Options: &definition.Index,
			}
			collection.Indexes().CreateOne(context.TODO(), indexModel)
		}
	}
}
