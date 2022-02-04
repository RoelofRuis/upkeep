package app

import (
	"errors"
	"time"
	"timesheet/infra"
	"timesheet/model/repo"
)

type Repository repo.Repository

func (r Repository) Purge() infra.Handler {
	return func(args []string) (error, string) {
		timesheet, err := r.Timesheets.GetForDay(time.Now())
		if err != nil {
			return err, ""
		}

		if err := r.Timesheets.Delete(timesheet); err != nil {
			return err, ""
		}

		return nil, "purged"
	}
}

func (r Repository) Edit(f func(args []string, editor TimesheetEditor) (error, string)) infra.Handler {
	return func(args []string) (error, string) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return err, ""

		}
		timesheet, err := r.Timesheets.GetForDay(time.Now())
		if err != nil {
			return err, ""
		}

		err, s := f(args, TimesheetEditor{upkeep: upkeep, timesheet: timesheet})
		if err != nil {
			return err, s
		}

		if err := r.Upkeep.Insert(upkeep); err != nil {
			return err, ""
		}
		if err := r.Timesheets.Insert(timesheet); err != nil {
			return err, ""
		}

		return nil, s
	}
}

func HandleStart(args []string, editor TimesheetEditor) (error, string) {
	editor.Start(args)

	return nil, editor.Show()
}

func HandleStop(args []string, editor TimesheetEditor) (error, string) {
	editor.Stop()

	return nil, editor.Show()
}

func HandleSwitch(args []string, editor TimesheetEditor) (error, string) {
	editor.Switch(args)

	return nil, editor.Show()
}

func HandleTag(args []string, editor TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no tags specified"), ""
	}

	editor.Tag(args)

	return nil, editor.Show()
}

func HandleShow(args []string, editor TimesheetEditor) (error, string) {
	return nil, editor.Show()
}
