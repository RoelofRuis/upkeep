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

func (r Repository) Edit(f func(args []string, editor *TimesheetEditor) (string, error)) infra.Handler {
	return func(args []string) (string, error) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return "", err
		}
		timesheet, err := r.Timesheets.GetForDate(model.NewDate(time.Now()))
		if err != nil {
			return "", err
		}

		editor := &TimesheetEditor{upkeep: &upkeep, timesheet: &timesheet}

		s, err := f(args, editor)
		if err != nil {
			return s, err
		}

		if editor.upkeep != &upkeep {
			if err := r.Upkeep.Insert(*editor.upkeep); err != nil {
				return "", err
			}
		}
		if editor.timesheet != &timesheet {
			if err := r.Timesheets.Insert(*editor.timesheet); err != nil {
				return "", err
			}
		}

		return s, nil
	}
}

func (t *TimesheetEditor) Start(category string) {
	t.Stop()

	now := time.Now()
	sheet := t.timesheet.Start(now)

	quotum := t.upkeep.GetQuotumForDay(now.Weekday())
	sheet = sheet.SetQuotum(quotum)

	t.timesheet = &sheet

	if category != "" {
		t.Category(category)
	}
}

func (t *TimesheetEditor) Stop() {
	sheet := t.timesheet.Stop(time.Now(), t.upkeep.GetSelectedCategory().Name)
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Abort() {
	sheet := t.timesheet.Abort()
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Switch(category string) {
	t.Stop()
	upkeep := t.upkeep.ShiftSelectedCategory()
	t.upkeep = &upkeep
	t.Start(category)
}

func (t *TimesheetEditor) Continue(category string) {
	t.Stop()
	upkeep := t.upkeep.UnshiftSelectedCategory()
	t.upkeep = &upkeep
	t.Start(category)
}

func (t *TimesheetEditor) Remove(blockId int) {
	timesheet := t.timesheet.RemoveBlock(blockId)
	t.timesheet = &timesheet
}

func (t *TimesheetEditor) Update(blockId int, category string) {
	timesheet := t.timesheet.UpdateBlockCategory(blockId, category)
	t.timesheet = &timesheet
}

var validCategory = regexp.MustCompile(`^[a-z0-9_]+$`)

func (t *TimesheetEditor) Category(category string) {
	upkeep := *t.upkeep

	if !validCategory.MatchString(category) {
		return
	}
	upkeep = upkeep.SetSelectedCategory(category)
	t.upkeep = &upkeep
}

func (t *TimesheetEditor) SetCategoryMaxDayQuotum(cat string, dur *time.Duration) {
	upkeep := t.upkeep.SetCategoryMaxDayQuotum(cat, dur)
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
	return ViewSheet(*t.upkeep, *t.timesheet)
}
