package e_tests

import (
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCoreSchedule(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

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

		time.Sleep(1 * time.Second)
	})
}
