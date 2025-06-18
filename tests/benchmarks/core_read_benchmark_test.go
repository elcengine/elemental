//go:build benchmark

package benchmarks

import (
	"testing"

	. "github.com/elcengine/elemental/tests/fixtures"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BenchmarkCoreRead(b *testing.B) {
	ts.SeededConnection(b.Name())

	UserModel := UserModel.SetDatabase(b.Name())

	for b.Loop() {
		UserModel.Find().ExecTT(b.Context())
	}
}

func BenchmarkCoreReadDriver(b *testing.B) {
	ts.SeededConnection(b.Name())

	coll := UserModel.SetDatabase(b.Name()).Collection()

	for b.Loop() {
		var users []User
		cursor := lo.Must(coll.Find(b.Context(), primitive.M{}))
		lo.Must0(cursor.All(b.Context(), &users))
	}
}
