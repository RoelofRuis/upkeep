package main

import (
	"strings"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

func ViewDays(req *Request) (string, error) {
	printer := infra.TerminalPrinter{}

	groupCategories := GroupCategories(req.Params)

	upToRecentDur := time.Duration(0)
	totalDur := time.Duration(0)
	upToRecentQuotum := model.NewDuration()
	totalQuotum := model.NewDuration()

	for _, daySheet := range req.Timesheets {
		blocks := req.Upkeep.DiscountTimeBlocks(*daySheet, req.Clock.Now())
		dayDur := blocks.TotalDuration()
		totalDur += dayDur

		dayQuotum := req.Upkeep.GetTimesheetQuotum(*daySheet)
		totalQuotum = totalQuotum.Add(dayQuotum)

		code := infra.Bold
		if daySheet.Finalised {
			code += infra.Green
		}
		if !daySheet.Date.After(req.Clock.Now()) {
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

	if upToRecentQuotum.IsZero() {
		printer.PrintC(infra.Bold, "                %s", infra.FormatDurationBracketed(upToRecentDur))
	} else {
		printer.PrintC(infra.Bold, "                %s / %s",
			infra.FormatDurationBracketed(upToRecentDur),
			infra.FormatDurationBracketed(upToRecentQuotum.Get()),
		).
			Print(" %s", infra.FormatPercentage(upToRecentDur, upToRecentQuotum.Get())).
			Newline()
	}

	if !totalQuotum.IsZero() {
		printer.
			PrintC(infra.White+infra.Bold, "                %s / %s",
				infra.FormatDurationBracketed(totalDur),
				infra.FormatDurationBracketed(totalQuotum.Get()),
			).
			Print(" %s", infra.FormatPercentage(totalDur, totalQuotum.Get()))
	}

	return printer.String(), nil
}
