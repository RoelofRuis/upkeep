package main

import (
	"strings"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

func ViewDays(app *App) (string, error) {
	printer := infra.TerminalPrinter{}

	groupCategories := GroupCategories(app.Params)

	upToRecentDur := time.Duration(0)
	totalDur := time.Duration(0)
	upToRecentQuotum := model.NewDuration()
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
		if !daySheet.Date.After(time.Now()) {
			upToRecentDur += dayDur
			upToRecentQuotum = upToRecentQuotum.Add(dayQuotum)
		} else {
			code += infra.White
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
			printer.PrintC(infra.Bold, "%s / %s ", infra.FormatDurationBracketed(dayDur), infra.FormatDurationBracketed(dayQuotum.Get()))
		}

		printer.PrintC(infra.Green, "%s", strings.Join(daySheet.GetCategoryNames(groupCategories), " ")).Newline()
	}

	upToRecentPerc := (float64(upToRecentDur) / float64(upToRecentQuotum.Get())) * 100

	printer.PrintC(infra.Bold, "                %s / %s (%0.1f%%)",
		infra.FormatDurationBracketed(upToRecentDur),
		infra.FormatDurationBracketed(upToRecentQuotum.Get()),
		upToRecentPerc,
	).Newline()

	totalPerc := (float64(totalDur) / float64(totalQuotum.Get())) * 100

	printer.PrintC(infra.White+infra.Bold, "                %s / %s (%0.1f%%)",
		infra.FormatDurationBracketed(totalDur),
		infra.FormatDurationBracketed(totalQuotum.Get()),
		totalPerc,
	)

	return printer.String(), nil
}
