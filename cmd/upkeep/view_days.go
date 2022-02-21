package main

import (
	"strings"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

func ViewDays(app *App) (string, error) {
	printer := infra.TerminalPrinter{}

	totalDur := time.Duration(0)
	totalQuotum := model.NewDuration()
	for _, daySheet := range app.Timesheets {
		blocks := app.Upkeep.DiscountTimeBlocks(*daySheet, time.Now())
		dayDur := blocks.TotalDuration()
		totalDur += dayDur

		dayQuotum := app.Upkeep.GetTimesheetQuotum(*daySheet)
		totalQuotum = totalQuotum.Add(dayQuotum)

		code := infra.Bold
		if daySheet.Finalised {
			code += infra.Green
		}
		printer.PrintC(code, "%s ", daySheet.Date.Format("Mon 02 Jan 2006"))

		if dayDur == 0 && !dayQuotum.IsDefined() {
			printer.Newline()
			continue
		}

		if !dayQuotum.IsDefined() {
			printer.PrintC(infra.Bold, "%s", infra.FormatDurationBracketed(dayDur)).
				Print("           ")
		} else {
			printer.PrintC(infra.Bold, "%s", infra.FormatDurationBracketed(dayDur)).
				Print(" / %s ", infra.FormatDurationBracketed(dayQuotum.Get()))
		}

		printer.PrintC(infra.Green, "%s", strings.Join(daySheet.GetCategoryNames(), " ")).Newline()
	}

	totalPerc := (float64(totalDur) / float64(totalQuotum.Get())) * 100

	printer.PrintC(infra.Bold, "                %s", infra.FormatDurationBracketed(totalDur)).
		Print(" / %s (%0.1f%%)", infra.FormatDurationBracketed(totalQuotum.Get()), totalPerc)

	return printer.String(), nil
}
