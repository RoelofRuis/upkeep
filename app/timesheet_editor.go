package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"timesheet/infra"
	"timesheet/model"
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
	now := time.Now()
	sheet := t.timesheet.Start(now)

	quotum := t.upkeep.GetQuotumForDay(now.Weekday())
	if quotum != 0 {
		sheet = sheet.SetQuotum(quotum)
	}

	t.timesheet = &sheet

	t.Tag(tags)
}

func (t *TimesheetEditor) Switch(tags []string) {
	t.Stop()
	t.Start(tags)
}

func (t *TimesheetEditor) Stop() {
	sheet := t.timesheet.Stop(time.Now(), t.upkeep.GetTags())
	t.timesheet = &sheet
}

func (t *TimesheetEditor) Abort() {
	sheet := t.timesheet.Abort()
	t.timesheet = &sheet
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
			"%s - %s [%s] <%s>",
			block.Start.Format(model.LayoutHour),
			block.End.Format(model.LayoutHour),
			formatDur(block.Duration()),
			block.Tags.String(),
		))
	}

	if t.timesheet.LastStart.IsStarted() {
		lines = append(lines, fmt.Sprintf(
			"%s -              <%s>",
			t.timesheet.LastStart.Format(model.LayoutHour),
			t.upkeep.GetTags().String(),
		))
	}

	quotum := t.timesheet.Quotum

	if quotum == 0 {
		lines = append(lines, fmt.Sprintf(
			"              [%s]",
			formatDur(t.timesheet.Duration()),
		))
	} else {
		perc := (float64(t.timesheet.Duration()) / float64(quotum)) * 100

		lines = append(lines, fmt.Sprintf(
			"              [%s] / [%s] (%0.2f%%)",
			formatDur(t.timesheet.Duration()),
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
