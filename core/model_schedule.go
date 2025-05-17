package elemental

import (
	robocron "github.com/robfig/cron/v3"
)

var cron = robocron.New(robocron.WithSeconds())

// Marks this query to be executed on a given schedule.
// For the schedule format, see https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format
func (m Model[T]) Schedule(spec string) Model[T] {
	m.schedule = &spec
	return m
}

// Unschedule a query that was previously scheduled.
func (m Model[T]) Unschedule(id int) {
	cron.Remove(robocron.EntryID(id))
	if len(cron.Entries()) == 0 {
		cron.Stop()
	}
}

// Unschedule all queries that were previously scheduled.
func (m Model[T]) UnscheduleAll() {
	entries := cron.Entries()
	for _, entry := range entries {
		cron.Remove(entry.ID)
	}
	cron.Stop()
}
