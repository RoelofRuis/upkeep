package app

import (
	"fmt"
	"regexp"
	"strings"
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

func (t *TimesheetEditor) Start(tags []string) {
	t.Stop()

	now := time.Now()
	sheet := t.timesheet.Start(now)

	quotum := t.upkeep.GetQuotumForDay(now.Weekday())
	if quotum != 0 {
		sheet = sheet.SetQuotum(quotum)
	}

	t.timesheet = &sheet

	if tags != nil {
		t.upkeep.ClearTags()
		t.Tag(tags)
	}
}

func (t *TimesheetEditor) Stop() {
	sheet := t.timesheet.Stop(time.Now(), t.upkeep.GetTags())
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Abort() {
	sheet := t.timesheet.Abort()
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Switch(tags []string) {
	t.Stop()
	upkeep := t.upkeep.ShiftTags()
	t.upkeep = &upkeep
	t.Start(tags)
}

func (t *TimesheetEditor) Continue() {
	t.Stop()
	upkeep := t.upkeep.UnshiftTags()
	t.upkeep = &upkeep
	t.Start(nil)
}

var validTag = regexp.MustCompile(`^[+-]?[a-z_]*$`)

func (t *TimesheetEditor) Tag(tags []string) {
	upkeep := *t.upkeep
	for _, tag := range tags {
		if !validTag.MatchString(tag) {
			continue
		}
		if strings.HasPrefix(tag, "-") {
			upkeep = upkeep.RemoveTag(strings.TrimPrefix(tag, "-"))
		} else {
			upkeep = upkeep.AddTag(strings.TrimPrefix(tag, "+"))
		}
	}
	t.upkeep = &upkeep
}

func (t *TimesheetEditor) Show() string {
	var lines []string
	lines = append(lines, fmt.Sprintf(
		"@ %s\n<%s>",
		t.timesheet.Created.Format("Monday 02 Jan 2006"),
		t.upkeep.Tags.String(),
	))

	for _, block := range t.timesheet.Blocks {
		lines = append(lines, fmt.Sprintf(
			"%2d [%s - %s] [%s] <%s>",
			block.Id,
			block.Start.Format(model.LayoutHour),
			block.End.Format(model.LayoutHour),
			formatDur(block.Duration()),
			block.Tags.String(),
		))
	}

	if t.timesheet.IsStarted() {
		start := t.timesheet.LastStart
		end := model.NewMoment().Start(time.Now())
		dur := end.Sub(start)

		lines = append(lines, fmt.Sprintf(
			"   [%s - %s) [%s] <%s>",
			start.Format(model.LayoutHour),
			end.Format(model.LayoutHour),
			formatDur(dur),
			t.upkeep.GetTags().String(),
		))
	}

	quotum := t.timesheet.Quotum
	totalDuration := t.upkeep.TimesheetDuration(*t.timesheet)

	if quotum == 0 {
		lines = append(lines, fmt.Sprintf(
			"                   [%s]",
			formatDur(totalDuration),
		))
	} else {
		perc := (float64(totalDuration) / float64(quotum)) * 100

		lines = append(lines, fmt.Sprintf(
			"                   [%s] / [%s] (%0.1f%%)",
			formatDur(totalDuration),
			formatDur(quotum),
			perc,
		))
	}

	return strings.Join(lines, "\n")
}

func formatDur(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) - (hours * 60)

	return fmt.Sprintf("%01d:%02d", hours, minutes)
}

func (t *TimesheetEditor) Purge() {
	sheet := model.NewTimesheet(time.Now())
	t.timesheet = &sheet
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
