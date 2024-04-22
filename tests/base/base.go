package e_test_base

import (
	"elemental/core"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Age        int                `json:"age" bson:"age"`
	Occupation string             `json:"occupation" bson:"occupation,omitempty"`
	Weapons    []string           `json:"weapons" bson:"weapons"`
	Retired    bool               `json:"retired" bson:"retired"`
	School     *string            `json:"school" bson:"school"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Castle struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Kingdom struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Monster struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Category  string             `json:"category" bson:"category"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

var DefaultAge = 18

var UserModel = elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     reflect.String,
		Required: true,
		Index: options.IndexOptions{
			Unique: lo.ToPtr(true),
		},
	},
	"Age": {
		Type:    reflect.Int,
		Default: DefaultAge,
	},
	"Occupation": {
		Type: reflect.String,
	},
	"Weapons": {
		Type:    reflect.Slice,
		Default: []string{},
	},
	"Retired": {
		Type:    reflect.Bool,
		Default: false,
	},
}, elemental.SchemaOptions{
	Collection: "users",
}))
