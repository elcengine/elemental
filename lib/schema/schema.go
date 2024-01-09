package elemental

import (
	"reflect"

	"github.com/creasty/defaults"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	Definitions map[string]Field
	Options     SchemaOptions
}

type Model struct {
	Name   string
	Schema Schema
}

func NewSchema(definitions map[string]Field, opts SchemaOptions) Schema {
	defaults.Set(opts)
	return Schema{
		Definitions: definitions,
		Options:     opts,
	}
}

// Enables timestamps with the default field names of createdAt and updatedAt.
//
// @returns void
//
// @example
//
// schema.DefaultTimestamps()
func (s Schema) DefaultTimestamps() {
	s.Timestamps(nil)
}

// Enables timestamps with custom field names.
//
// @param *ts* - A struct containing the custom field names.
//
// @returns void
//
// @example
//
//	schema.Timestamps(&TS{
//		CreatedAt: "created_at",
//		UpdatedAt: "updated_at",
//	})
func (s Schema) Timestamps(ts *TS) {
	defaults.Set(s.Options.Timestamps)
	s.Options.Timestamps.Enabled = true
	if ts.CreatedAt != "" {
		s.Options.Timestamps.CreatedAt = ts.CreatedAt
	}
	if ts.UpdatedAt != "" {
		s.Options.Timestamps.UpdatedAt = ts.UpdatedAt
	}
}

func (m Model) Validate() error {
	return nil
}

func (m Model) ValidateField() error {
	return nil
}

func (m Model) Create() (primitive.ObjectID, error) {
	return primitive.ObjectID{}, nil
}

func (m Model) FindOne(query primitive.M) *T {
	userSchema := NewSchema(map[string]Field{
		"ID": Field{
			Disabled: true,
		},
		"Name": Field{
			Type:     reflect.String,
			Required: true,
		},
	}, SchemaOptions{
		Collection: "users",
	})
	type User struct {
		Model
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
	user := User{
		Name: "John Doe",
	}
	u, _ := User.Create(user)
	User.FindOne(User{}, primitive.M{"name": "John Doe"})
}

type UserSchema struct {
	ID   bool `default:"true"`
	Name Field
}