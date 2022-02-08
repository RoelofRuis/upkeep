package app

import (
	"regexp"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

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

type TimesheetEditor struct {
	upkeep    *model.Upkeep
	timesheet *model.Timesheet
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

func (t *TimesheetEditor) Day() string {
	excludedCategories := t.upkeep.ExcludedCategories

	printer := infra.TerminalPrinter{}
	printer.Print("@ %s", t.timesheet.Created.Format("Monday 02 Jan 2006")).Newline()
	printer.Green("%s", t.upkeep.Categories.String()).Newline()

	for _, block := range t.timesheet.Blocks {
		printer.White("%2d ", block.Id).
			Print("[%s - %s]", block.Start.Format(model.LayoutHour), block.End.Format(model.LayoutHour))

		if excludedCategories.Contains(block.Category) {
			printer.Print(" [%s] ", infra.FormatDuration(block.Duration())).
				Yellow("%s", block.Category)
		} else {
			printer.Bold(" [%s] ", infra.FormatDuration(block.Duration())).
				Green("%s", block.Category)
		}

		printer.Newline()
	}

	if t.timesheet.IsStarted() {
		start := t.timesheet.LastStart
		end := model.NewMoment().Start(time.Now())
		dur := end.Sub(start)

		printer.White(">> ").
			Print("[%s - %s] ", start.Format(model.LayoutHour), end.Format(model.LayoutHour))

		if excludedCategories.Contains(t.upkeep.GetCategory()) {
			printer.Print("[%s]", infra.FormatDuration(dur)).
				Yellow(" %s", t.upkeep.GetCategory())
		} else {
			printer.Bold("[%s]", infra.FormatDuration(dur)).
				Green(" %s", t.upkeep.GetCategory())
		}

		printer.Newline()
	}

	quotum := t.timesheet.Quotum
	totalDuration := t.upkeep.TimesheetDuration(*t.timesheet)

	if quotum == 0 {
		printer.Print("                   ").
			Bold("[%s]", infra.FormatDuration(totalDuration)).
			Newline()
	} else {
		perc := (float64(totalDuration) / float64(quotum)) * 100

		printer.Print("                   ").
			Bold("[%s]", infra.FormatDuration(totalDuration)).
			Print(" / [%s] (%0.1f%%)", infra.FormatDuration(quotum), perc).
			Newline()
	}

	return printer.String()
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
