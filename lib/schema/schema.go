package elemental

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
	"elemental/lib"
)

type Schema[T any] struct {
	T,
	CreatedAt *time.Time
	UpdatedAt *time.Time
	options   SchemaOptions
}

// Enables timestamps with the default field names of createdAt and updatedAt.
//
// @returns void
//
// @example
// 
// schema.DefaultTimestamps()
func (s Schema[T]) DefaultTimestamps() {
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
// schema.Timestamps(&TS{
// 	CreatedAt: "created_at",
// 	UpdatedAt: "updated_at",
// })
func (s Schema[T]) Timestamps(ts *TS) {
	elemental.SetDefaults(s.options.timestamps)
	s.options.timestamps.enabled = true
	if ts.CreatedAt != "" {
		s.options.timestamps.createdAt = ts.CreatedAt
	}
	if ts.UpdatedAt != "" {
		s.options.timestamps.updatedAt = ts.UpdatedAt
	}
}

func (s Schema[T]) Validate() error {
	return nil
}

func (s Schema[T]) ValidateField() error {
	return nil
}

func (s Schema[T]) Create() (primitive.ObjectID, error) {

}

func (s Schema[T]) FindOne(query primitive.M) *T {

}

type User struct {
	// ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	// Email            string             `json:"email" bson:"email,omitempty"`
	// Password         string             `json:"password" bson:"password,omitempty"`
	// Organizations    []string           `json:"organizations" bson:"organizations"`
	// Verified         bool               `json:"verified" bson:"verified"`
	// VerificationCode *string            `json:"verification_code" bson:"verification_code,omitempty"`
	// Role             UserRole           `json:"role" bson:"role,omitempty"`
	// CreatedAt        string             `json:"created_at" bson:"created_at,omitempty"`
	// UpdatedAt        string             `json:"updated_at" bson:"updated_at,omitempty"`
	Name Field
}

func (u User) Validate() error {
	a := User{
		Name: Field{
			Type: reflect.String,
		},
	}
	userSchema := Schema[User]{}
	userSchema.Timestamps(&TS{
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
	})
	fmt.Println(a)
	return nil
}
