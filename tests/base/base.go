package e_test_base

import (
	"elemental/core"
	"reflect"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Age        int                `json:"age" bson:"age"`
	Occupation string             `json:"occupation" bson:"occupation,omitempty"`
	Weapons    []string           `json:"weapons" bson:"weapons"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

var DefaultAge = 18

var UserModel = elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     reflect.String,
		Required: true,
		Index: options.IndexOptions{
			Unique: qkit.ValPtr(true),
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
}, elemental.SchemaOptions{
	Collection: "users",
}))
