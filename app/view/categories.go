package view

import (
	"fmt"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func ViewCategories(upkeep model.Upkeep, sheets []model.Timesheet) string {
	categoryDurations := make(map[string]time.Duration)
	nameLength := 0

	for _, sheet := range sheets {
		blocks := upkeep.DiscountTimeBlocks(sheet, time.Now())
		for _, block := range blocks {
			category := block.Block.Category
			nameLength = infra.Max(len(category), nameLength)

			dur, has := categoryDurations[category]
			if !has {
				dur = time.Duration(0)
			}
			dur += block.DiscountedDuration
			categoryDurations[category] = dur
		}
	}

	printer := infra.TerminalPrinter{}

	for cat, dur := range categoryDurations {
		format := fmt.Sprintf("%%-%ds", nameLength)
		printer.Green(format, cat).
			Print(" %s", infra.FormatDuration(dur)).
			Newline()
	}

	return printer.String()
}
