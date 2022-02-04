package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"timesheet/model"
)

type TimesheetEditor struct {
	upkeep    *model.Upkeep
	timesheet *model.Timesheet
}

func (t *TimesheetEditor) Start(tags []string) {
	t.timesheet.Start(time.Now())
	t.Tag(tags)
}

func (t *TimesheetEditor) Switch(tags []string) {
	t.Stop()
	t.Start(tags)
}

func (t *TimesheetEditor) Stop() {
	t.timesheet.Stop(time.Now(), t.upkeep.GetTags())
}

func (t *TimesheetEditor) Abort() {
	t.timesheet.Abort()
}

var validTag = regexp.MustCompile(`^[+-]?[a-z_]*$`)

func (t *TimesheetEditor) Tag(tags []string) {
	for _, tag := range tags {
		if !validTag.MatchString(tag) {
			continue
		}
		if strings.HasPrefix(tag, "-") {
			t.upkeep.RemoveTag(strings.TrimPrefix(tag, "-"))
		} else {
			t.upkeep.AddTag(strings.TrimPrefix(tag, "+"))
		}
	}
}

func (t *TimesheetEditor) Show() string {
	var lines []string
	lines = append(lines, fmt.Sprintf(
		"> %s\ntags [%s]",
		t.timesheet.Created.Format("Monday 02 Jan 2006"),
		t.upkeep.Tags.String(),
	))

	for _, block := range t.timesheet.Blocks {
		lines = append(lines, fmt.Sprintf(
			"%s - %s (%s) [%s]",
			block.Start.Format(model.LayoutHour),
			block.End.Format(model.LayoutHour),
			formatDur(block.Duration()),
			block.Tags.String(),
		))
	}

	if t.timesheet.LastStart.IsStarted() {
		lines = append(lines, fmt.Sprintf(
			"%s -              [%s]",
			t.timesheet.LastStart.Format(model.LayoutHour),
			t.upkeep.GetTags().String(),
		))
	}

	quotum := t.upkeep.QuotumForDay(t.timesheet.Created.Weekday())

	if quotum == 0 {
		lines = append(lines, fmt.Sprintf(
			"              (%s)",
			formatDur(t.timesheet.Duration()),
		))
	} else {
		lines = append(lines, fmt.Sprintf(
			"              (%s) of (%s)",
			formatDur(t.timesheet.Duration()),
			formatDur(quotum),
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
	t.timesheet = model.NewTimesheet(time.Now())
}
