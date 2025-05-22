// Shared models and types for all test suites.
package fixtures

import (
	"time"

	elemental "github.com/elcengine/elemental/core"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type MonsterWeakness struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Oils           []string           `json:"oils" bson:"oils"`
	Signs          []string           `json:"signs" bson:"signs"`
	Decoctions     []string           `json:"decoctions" bson:"decoctions"`
	Bombs          []string           `json:"bombs" bson:"bombs"`
	InvulnerableTo []string           `json:"invulnerable_to" bson:"invulnerable_to"`
}

type Monster struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Category   string             `json:"category,omitempty" bson:"category,omitempty"`
	Weaknesses MonsterWeakness    `json:"weaknesses" bson:"weaknesses"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Bestiary struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Monster Monster            `json:"monster" bson:"monster"`
	Kingdom Kingdom            `json:"kingdom" bson:"kingdom"`
}

type BestiaryWithID struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	MonsterID string             `json:"monster_id" bson:"monster_id"`
}

const DefaultUserAge = 18

var UserModel = elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     elemental.String,
		Required: true,
		Index:    options.Index().SetUnique(true),
	},
	"Age": {
		Type:    elemental.Int,
		Default: DefaultUserAge,
	},
	"Occupation": {
		Type: elemental.String,
	},
	"Weapons": {
		Type:    elemental.Slice,
		Default: []string{},
	},
	"Retired": {
		Type:    elemental.Bool,
		Default: false,
	},
}, elemental.SchemaOptions{
	Collection: "users",
}))

var MonsterModel = elemental.NewModel[Monster]("Monster", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     elemental.String,
		Required: true,
	},
	"Category": {
		Type: elemental.String,
	},
	"Weaknesses": {
		Type: elemental.Struct,
		Schema: lo.ToPtr(elemental.NewSchema(map[string]elemental.Field{
			"Oils": {
				Type: elemental.Slice,
			},
			"Signs": {
				Type:    elemental.Slice,
				Default: []string{"Igni"},
			},
			"Decoctions": {
				Type: elemental.Slice,
			},
			"Bombs": {
				Type: elemental.Slice,
			},
			"InvulnerableTo": {
				Type:    elemental.Slice,
				Default: []string{"Steel"},
			},
		})),
	},
}, elemental.SchemaOptions{
	Collection: "monsters",
}))

var KingdomModel = elemental.NewModel[Kingdom]("Kingdom", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     elemental.String,
		Required: true,
	},
}, elemental.SchemaOptions{
	Collection: "kingdoms",
}))

var BestiaryModel = elemental.NewModel[Bestiary]("Bestiary", elemental.NewSchema(map[string]elemental.Field{
	"Monster": {
		Type: elemental.Struct,
		Ref:  "Monster",
	},
	"Kingdom": {
		Type: elemental.ObjectID,
		Ref:  "Kingdom",
	},
}, elemental.SchemaOptions{
	Collection: "bestiary",
}))

var BestiaryWithIDModel = elemental.NewModel[BestiaryWithID]("BestiaryWithID", elemental.NewSchema(map[string]elemental.Field{
	"MonsterID": {
		Type:    elemental.String,
		Ref:     "Monster",
		IsRefID: true,
	},
}, elemental.SchemaOptions{
	Collection: "bestiaryWithID",
}))
