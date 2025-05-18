package tests

import (
	"testing"
	"time"

	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreSchedule(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	KingdomModel := KingdomModel.SetDatabase(t.Name())

	Convey("Schedule document creation every 2 seconds", t, func() {
		id := KingdomModel.Create(Kingdom{
			Name: uuid.NewString(),
		}).Schedule("*/2 * * * * *").ExecInt()

		defer KingdomModel.Unschedule(id)

		for i := range 3 {
			SoTimeout(t, func() (ok bool) {
				if len(KingdomModel.Find().ExecTT()) >= i {
					ok = true
				}
				return
			})
			time.Sleep(2 * time.Second)
		}
	})
	Convey("Schedule document creation which should error out", t, func() {
		executionErrorOcurred := false
		id := KingdomModel.SetConnection("I don't exist").Create(Kingdom{
			Name: uuid.NewString(),
		}).Schedule("*/2 * * * * *", func(a any) {
			executionErrorOcurred = true
		}).ExecInt()

		defer KingdomModel.Unschedule(id)

		SoTimeout(t, func() (ok bool) {
			if executionErrorOcurred {
				ok = true
			}
			return
		})
	})
}
