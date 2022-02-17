package view

import (
	"strings"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func ViewSheets(upkeep model.Upkeep, timesheets []model.Timesheet) (string, error) {
	printer := infra.TerminalPrinter{}
	printer.PrintC(infra.BGGreen, "%s", strings.Join(upkeep.SelectedCategories, " | ")).Newline()

	for _, timesheet := range timesheets {
		printer.PrintC(infra.Bold,"@ %s", timesheet.Date.Format("Monday 02 Jan 2006")).Newline()

		blocks := upkeep.DiscountTimeBlocks(timesheet, time.Now())

		for _, block := range blocks {
			if block.Block.Id == -1 {
				printer.PrintC(infra.White, ">> ")
			} else {
				printer.PrintC(infra.White, "%2d ", block.Block.Id)
			}

			if block.Block.HasEnded() {
				if block.Block.Type == model.TypeTime {
					printer.Print("[%s - %s]",
						block.Block.WithTime.Start.Format(model.LayoutHour),
						block.Block.WithTime.End.Format(model.LayoutHour),
					)
				} else {
					printer.Print("               ")
				}

				if block.IsDiscounted {
					if block.DiscountedDuration != 0 {
						printer.PrintC(infra.Yellow, " %s ", infra.FormatDurationBracketed(block.DiscountedDuration))
					} else {
						printer.Print("         ")
					}
				} else {
					printer.PrintC(infra.Bold," %s ", infra.FormatDurationBracketed(block.Block.BaseDuration()))
				}
			} else {
				printer.PrintC(infra.Red, "[%s -   ?  ]         ",
					block.Block.WithTime.Start.Format(model.LayoutHour),
				)
			}

			printer.PrintC(infra.Green, "%s", block.Block.Category)

			printer.Newline()
		}

		quotum := upkeep.GetTimesheetQuotum(timesheet)
		totalDuration := blocks.TotalDuration()

		if !quotum.IsDefined() {
			printer.Print("                   ").
				PrintC(infra.Bold,"%s", infra.FormatDurationBracketed(totalDuration)).
				Newline()
		} else {
			perc := (float64(totalDuration) / float64(quotum.Get())) * 100

			printer.Print("                   ").
				PrintC(infra.Bold, "%s", infra.FormatDurationBracketed(totalDuration)).
				Print(" / %s (%0.1f%%)", infra.FormatDurationBracketed(quotum.Get()), perc).
				Newline()
		}
	}

	return printer.String(), nil
}
