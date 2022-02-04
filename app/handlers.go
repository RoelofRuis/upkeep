package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"timesheet/infra"
	"timesheet/model/repo"
)

type Repository repo.Repository

func (r Repository) Edit(f func(args []string, editor *TimesheetEditor) (error, string)) infra.Handler {
	return func(args []string) (error, string) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return err, ""

		}
		timesheet, err := r.Timesheets.GetForDay(time.Now())
		if err != nil {
			return err, ""
		}

		editor := &TimesheetEditor{upkeep: upkeep, timesheet: timesheet}

		err, s := f(args, editor)
		if err != nil {
			return err, s
		}

		if err := r.Upkeep.Insert(editor.upkeep); err != nil {
			return err, ""
		}
		if err := r.Timesheets.Insert(editor.timesheet); err != nil {
			return err, ""
		}

		return nil, s
	}
}

func HandlePurge(args []string, editor *TimesheetEditor) (error, string) {
	editor.Purge()

	return nil, editor.Show()
}

func HandleStart(args []string, editor *TimesheetEditor) (error, string) {
	editor.Start(args)

	return nil, editor.Show()
}

func HandleStop(args []string, editor *TimesheetEditor) (error, string) {
	editor.Stop()

	return nil, editor.Show()
}

func HandleAbort(args []string, editor *TimesheetEditor) (error, string) {
	editor.Abort()

	return nil, editor.Show()
}

func HandleSwitch(args []string, editor *TimesheetEditor) (error, string) {
	editor.Switch(args)

	return nil, editor.Show()
}

func HandleTag(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no tags specified"), ""
	}

	editor.Tag(args)

	return nil, editor.Show()
}

func HandleShow(args []string, editor *TimesheetEditor) (error, string) {
	return nil, editor.Show()
}

func HandleConf(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no setting given"), ""
	}

	if args[0] == "quotum" {
		if len(args) < 3 {
			return errors.New("too few arguments"), ""
		}
		weekday, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err, ""
		}
		duration, err := time.ParseDuration(args[2])
		if err != nil {
			return err, ""
		}
		editor.SetQuotum(time.Weekday(weekday), duration)
		return nil, fmt.Sprintf("updated '%s'", args[0])
	}

	return fmt.Errorf("unknown setting '%s'", args[0]), ""
}