package view

import (
	"strings"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func ViewDays(upkeep model.Upkeep, timesheets []model.Timesheet) (string, error) {
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
			printer.PrintC(infra.Bold,"%s", infra.FormatDurationBracketed(dayDur))
		} else {
			printer.PrintC(infra.Bold, "%s", infra.FormatDurationBracketed(dayDur)).
				Print(" / %s ", infra.FormatDurationBracketed(dayQuotum.Get()))
		}

		printer.PrintC(infra.Green, " %s", strings.Join(daySheet.GetCategoryNames(), " ")).Newline()
	}

	totalPerc := (float64(totalDur) / float64(totalQuotum.Get())) * 100

	printer.PrintC(infra.Bold,"                %s", infra.FormatDurationBracketed(totalDur)).
		Print(" / %s (%0.1f%%)", infra.FormatDurationBracketed(totalQuotum.Get()), totalPerc)

	return printer.String(), nil
}
