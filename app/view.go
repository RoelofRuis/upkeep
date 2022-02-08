package app

import (
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func (r Repository) HandleViewAt(time time.Time) infra.Handler {
	return func(args []string) (error, string) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return err, ""
		}
		timesheet, err := r.Timesheets.GetForDay(time)
		if err != nil {
			return err, ""
		}

		return nil, ViewDay(upkeep, timesheet)
	}
}

func ViewDay(upkeep model.Upkeep, timesheet model.Timesheet) string {
	excludedCategories := upkeep.ExcludedCategories

	printer := infra.TerminalPrinter{}
	printer.Print("@ %s", timesheet.Created.Format("Monday 02 Jan 2006")).Newline()
	printer.Green("%s", upkeep.Categories.String()).Newline()

	for _, block := range timesheet.Blocks {
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

	if timesheet.IsStarted() {
		start := timesheet.LastStart
		end := model.NewMoment().Start(time.Now())
		dur := end.Sub(start)

		printer.White(">> ").
			Print("[%s - %s] ", start.Format(model.LayoutHour), end.Format(model.LayoutHour))

		if excludedCategories.Contains(upkeep.GetCategory()) {
			printer.Print("[%s]", infra.FormatDuration(dur)).
				Yellow(" %s", upkeep.GetCategory())
		} else {
			printer.Bold("[%s]", infra.FormatDuration(dur)).
				Green(" %s", upkeep.GetCategory())
		}

		printer.Newline()
	}

	quotum := timesheet.Quotum
	totalDuration := upkeep.TimesheetDuration(timesheet)

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
