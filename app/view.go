package app

import (
	"fmt"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func (r Repository) HandleViewWeek(args []string) (error, string) {
	date := model.Today().PreviousMonday()

	upkeep, err := r.Upkeep.Get()
	if err != nil {
		return err, ""
	}

	sheets := make([]model.Timesheet, 5)
	for i, day := range date.Week(5) {
		sheet, err := r.Timesheets.GetForDate(day)
		if err != nil {
			return err, ""
		}
		sheets[i] = sheet
	}

	return nil, ViewWeek(upkeep, sheets)
}

func ViewWeek(upkeep model.Upkeep, sheets []model.Timesheet) string {
	printer := infra.TerminalPrinter{}

	weekDur := time.Duration(0)
	weekQuotum := time.Duration(0)
	for _, daySheet := range sheets {
		dayDur := upkeep.TimesheetDuration(daySheet)
		weekDur += dayDur

		dayQuotum := upkeep.TimesheetQuotum(daySheet)
		weekQuotum += dayQuotum

		printer.Print("%s ", daySheet.Date.Format("Mon 02 Jan 2006"))

		if dayDur == 0 && dayQuotum == 0 {
			printer.Newline()
			continue;
		}

		if dayQuotum == 0 {
				printer.Bold("[%s]", infra.FormatDuration(dayDur))
		} else {
			printer.Bold("[%s]", infra.FormatDuration(dayDur)).
				Print(" / [%s] ", infra.FormatDuration(dayQuotum))
		}

		printer.Green("%s", daySheet.GetCategories().String()).Newline()
	}

	weekPerc := (float64(weekDur) / float64(weekQuotum)) * 100

	printer.Bold("                [%s]", infra.FormatDuration(weekDur)).
		Print(" / [%s] (%0.1f%%)", infra.FormatDuration(weekQuotum), weekPerc)

	return printer.String()
}

func (r Repository) HandleViewSheet(args []string) (error, string) {
	date := model.Today()
	if len(args) > 0 {
		switch args[0] {
		case "today":
			break
		case "yesterday":
			date = date.ShiftDay(-1)
			break
		default:
			parsedDate, err := model.NewDateFromString(args[0])
			if err != nil {
				return fmt.Errorf("invalid date value '%s'", args[0]), ""
			}
			date = parsedDate
			break
		}
	}

	upkeep, err := r.Upkeep.Get()
	if err != nil {
		return err, ""
	}
	timesheet, err := r.Timesheets.GetForDate(date)
	if err != nil {
		return err, ""
	}

	return nil, ViewSheet(upkeep, timesheet)
}

func ViewSheet(upkeep model.Upkeep, timesheet model.Timesheet) string {
	excludedCategories := upkeep.ExcludedCategories

	printer := infra.TerminalPrinter{}
	printer.Print("@ %s", timesheet.Date.Format("Monday 02 Jan 2006")).Newline()
	printer.BGGreen("%s", upkeep.Categories.String()).Newline()

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
		if timesheet.Date.IsToday() {
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
		} else {
			printer.Red(">> [%s -   ?  ]", start.Format(model.LayoutHour)).Newline()
		}
	}

	quotum := upkeep.TimesheetQuotum(timesheet)
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
