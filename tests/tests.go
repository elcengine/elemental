package tests

import (
	"testing"
	"time"

	"github.com/elcengine/elemental/tests/fixtures"
	. "github.com/smartystreets/goconvey/convey"
)

type Castle = fixtures.Castle

type Kingdom = fixtures.Kingdom

type Monster = fixtures.Monster

type Bestiary = fixtures.Bestiary

type BestiaryWithID = fixtures.BestiaryWithID

type MonsterWeakness = fixtures.MonsterWeakness

type User = fixtures.User

var UserModel = fixtures.UserModel

var MonsterModel = fixtures.MonsterModel

var KingdomModel = fixtures.KingdomModel

var BestiaryModel = fixtures.BestiaryModel

var BestiaryWithIDModel = fixtures.BestiaryWithIDModel

// Test helper function to wait for a condition to be true or timeout.
// It will keep checking the condition every 100 milliseconds until the timeout is reached.
// If the condition is true, the assertion will pass.
func SoTimeout(t *testing.T, f func() bool, timeout ...<-chan time.Time) {
	t.Helper()
	if len(timeout) == 0 {
		timeout = append(timeout, time.After(10*time.Second))
	}
	for {
		select {
		case <-timeout[0]:
			t.Errorf("Timeout waiting for assertion to execute")
		default:
			ok := f()
			if ok {
				So(ok, ShouldBeTrue)
				return
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
