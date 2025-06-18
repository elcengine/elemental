//go:build benchmark

package benchmarks

import (
	"testing"
	"time"

	. "github.com/elcengine/elemental/tests/fixtures"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BenchmarkCoreCreate(b *testing.B) {
	ts.Connection(b.Name())

	UserModel := UserModel.SetDatabase(b.Name())

	for b.Loop() {
		UserModel.Create(User{
			Name: uuid.NewString(),
		}).ExecT()
	}
}

func BenchmarkCoreCreateDriver(b *testing.B) {
	ts.Connection(b.Name())

	coll := UserModel.SetDatabase(b.Name()).Collection()

	for b.Loop() {
		coll.InsertOne(b.Context(), User{
			Name: uuid.NewString(),
		})
	}
}

func BenchmarkCoreCreateMany(b *testing.B) {
	ts.Connection(b.Name())

	UserModel := UserModel.SetDatabase(b.Name())

	docs := lo.Map(mocks.Users, func(user User, _ int) User {
		user.Name = uuid.NewString()
		return user
	})

	for b.Loop() {
		UserModel.CreateMany(docs).ExecTT()
		UserModel.Drop()
	}
}

func BenchmarkCoreCreateManyDriver(b *testing.B) {
	ts.Connection(b.Name())

	coll := UserModel.SetDatabase(b.Name()).Collection()

	docs := lo.Map(mocks.Users, func(user User, _ int) any {
		user.ID = primitive.NewObjectID()
		user.Name = uuid.NewString()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		return user
	})

	for b.Loop() {
		coll.InsertMany(b.Context(), docs)
		coll.Drop(b.Context())
	}
}
