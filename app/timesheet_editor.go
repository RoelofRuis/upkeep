package app

import (
	"regexp"
	"time"
	"upkeep/app/view"
	"upkeep/infra"
	"upkeep/model"
)

type TimesheetEditor struct {
	upkeep    *model.Upkeep
	timesheet *model.Timesheet
}

func (r Repository) Edit(f func(params infra.Params, editor *TimesheetEditor) (string, error)) infra.Handler {
	return func(params infra.Params) (string, error) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return "", err
		}

		dateParam, err := params.GetNamed("d")
		if err != nil {
			return "", err
		}

		date, _, err := MakeDateRange(model.NewDate(time.Now()), dateParam)
		if err != nil {
			return "", err
		}

		timesheet, err := r.Timesheets.GetForDate(date)
		if err != nil {
			return "", err
		}

		editor := &TimesheetEditor{upkeep: &upkeep, timesheet: &timesheet}

		s, err := f(params, editor)
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

func (r Repository) Read(view func(model.Upkeep, []model.Timesheet) (string, error)) func(params infra.Params) (string, error) {
	return func(params infra.Params) (string, error) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return "", err
		}

		dateParam, err := params.GetNamed("d")
		if err != nil {
			return "", err
		}

		date, numDays, err := MakeDateRange(model.NewDate(time.Now()), dateParam)
		if err != nil {
			return "", err
		}

		dates := date.IterateNext(numDays)
		timesheets := make([]model.Timesheet, len(dates))
		for i, day := range dates {
			sheet, err := r.Timesheets.GetForDate(day)
			if err != nil {
				return "", err
			}
			timesheets[i] = sheet
		}

		return view(upkeep, timesheets)
	}
}

func (t *TimesheetEditor) Start(category string) {
	t.Stop()

	now := time.Now()
	sheet := t.timesheet.Start(now)

	quotum := t.upkeep.GetWeekdayQuotum(now.Weekday())
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

func (t *TimesheetEditor) Restore(blockId int) {
	timesheet := t.timesheet.RestoreBlock(blockId)
	t.timesheet = &timesheet
}

func (t *TimesheetEditor) Update(blockId int, category string) {
	timesheet := t.timesheet.UpdateBlockCategory(blockId, category)
	t.timesheet = &timesheet
}

func (t *TimesheetEditor) Write(cat string, dur model.Duration) {
	timesheet := t.timesheet.Write(cat, dur)
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

func (t *TimesheetEditor) View() (string, error) {
	return view.ViewSheets(*t.upkeep, []model.Timesheet{*t.timesheet})
}
