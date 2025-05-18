//go:build benchmark

package benchmarks

import (
	"testing"

	. "github.com/elcengine/elemental/tests/fixtures"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
)

func BenchmarkCoreCreate(b *testing.B) {
	b.ReportAllocs()

	ts.Connection(b.Name())

	UserModel := UserModel.SetDatabase(b.Name())

	for b.Loop() {
		UserModel.Create(User{
			Name: uuid.NewString(),
		}).ExecT()
	}
}

func BenchmarkCoreCreateDriver(b *testing.B) {
	b.ReportAllocs()

	ts.Connection(b.Name())

	coll := UserModel.SetDatabase(b.Name()).Collection()

	for b.Loop() {
		coll.InsertOne(b.Context(), User{
			Name: uuid.NewString(),
		})
	}
}
