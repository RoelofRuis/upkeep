package view

import (
	"strings"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func ViewDays(upkeep model.Upkeep, timesheets []model.Timesheet) string {
	printer := infra.TerminalPrinter{}

	totalDur := time.Duration(0)
	totalQuotum := model.NewDuration()
	for _, daySheet := range timesheets {
		blocks := upkeep.DiscountTimeBlocks(daySheet, time.Now())
		dayDur := blocks.TotalDuration()
		totalDur += dayDur

		dayQuotum := upkeep.GetTimesheetQuotum(daySheet)
		totalQuotum = totalQuotum.Add(dayQuotum)

		printer.Print("%s ", daySheet.Date.Format("Mon 02 Jan 2006"))

		if dayDur == 0 && !dayQuotum.IsDefined() {
			printer.Newline()
			continue
		}

		if !dayQuotum.IsDefined() {
			printer.Bold("[%s]", infra.FormatDuration(dayDur))
		} else {
			printer.Bold("[%s]", infra.FormatDuration(dayDur)).
				Print(" / [%s] ", infra.FormatDuration(dayQuotum.Get()))
		}

		printer.Green(" %s", strings.Join(daySheet.GetCategoryNames(), " ")).Newline()
	}

	totalPerc := (float64(totalDur) / float64(totalQuotum.Get())) * 100

	printer.Bold("                [%s]", infra.FormatDuration(totalDur)).
		Print(" / [%s] (%0.1f%%)", infra.FormatDuration(totalQuotum.Get()), totalPerc)

	return printer.String()
}
