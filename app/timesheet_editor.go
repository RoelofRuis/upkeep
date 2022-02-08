package app

import (
	"regexp"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

type TimesheetEditor struct {
	upkeep    *model.Upkeep
	timesheet *model.Timesheet
}

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

		editor := &TimesheetEditor{upkeep: &upkeep, timesheet: &timesheet}

		err, s := f(args, editor)
		if err != nil {
			return err, s
		}

		if editor.upkeep != &upkeep {
			if err := r.Upkeep.Insert(*editor.upkeep); err != nil {
				return err, ""
			}
		}
		if editor.timesheet != &timesheet {
			if err := r.Timesheets.Insert(*editor.timesheet); err != nil {
				return err, ""
			}
		}

		return nil, s
	}
}

func (t *TimesheetEditor) Start(category string) {
	t.Stop()

	now := time.Now()
	sheet := t.timesheet.Start(now)

	quotum := t.upkeep.GetQuotumForDay(now.Weekday())
	if quotum != 0 {
		sheet = sheet.SetQuotum(quotum)
	}

	t.timesheet = &sheet

	if category != "" {
		t.Category(category)
	}
}

func (t *TimesheetEditor) Stop() {
	sheet := t.timesheet.Stop(time.Now(), t.upkeep.GetCategory())
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Abort() {
	sheet := t.timesheet.Abort()
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Switch(category string) {
	t.Stop()
	upkeep := t.upkeep.ShiftCategory()
	t.upkeep = &upkeep
	t.Start(category)
}

func (t *TimesheetEditor) Continue() {
	t.Stop()
	upkeep := t.upkeep.UnshiftCategory()
	t.upkeep = &upkeep
	t.Start("")
}

func (t *TimesheetEditor) Remove(blockId int) {
	timesheet := t.timesheet.Remove(blockId)
	t.timesheet = &timesheet
}

var validCategory = regexp.MustCompile(`^[a-z0-9_]+$`)

func (t *TimesheetEditor) Category(category string) {
	upkeep := *t.upkeep

	if !validCategory.MatchString(category) {
		return
	}
	upkeep = upkeep.SetCategory(category)
	t.upkeep = &upkeep
}

func (t *TimesheetEditor) Exclude(category string) {
	upkeep := t.upkeep.AddExcludedCategory(category)
	t.upkeep = &upkeep
}

func (t *TimesheetEditor) Include(category string) {
	upkeep := t.upkeep.RemoveExcludedCategory(category)
	t.upkeep = &upkeep
}

func (t *TimesheetEditor) AdjustQuotum(day time.Weekday, dur *time.Duration) {
	if dur == nil {
		upkeep := t.upkeep.RemoveQuotumForDay(day)
		t.upkeep = &upkeep
	} else {
		upkeep := t.upkeep.SetQuotumForDay(day, *dur)
		t.upkeep = &upkeep
	}
}

func (t *TimesheetEditor) View() string {
	return ViewDay(*t.upkeep, *t.timesheet)
}
