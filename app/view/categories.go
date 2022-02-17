package view

import (
	"fmt"
	"sort"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func ViewCategories(upkeep model.Upkeep, sheets []model.Timesheet) (string, error) {
	var categories []string
	durations := make(map[string]time.Duration)
	nameLength := 0

	for _, sheet := range sheets {
		blocks := upkeep.DiscountTimeBlocks(sheet, time.Now())
		for _, block := range blocks {
			category := block.Block.Category
			nameLength = infra.Max(len(category), nameLength)

			dur, has := durations[category]
			if !has {
				dur = time.Duration(0)
				categories = append(categories, category)
			}
			dur += block.DiscountedDuration
			durations[category] = dur
		}
	}

	sort.Slice(categories, func(i, j int) bool {
		return durations[categories[i]] > durations[categories[j]]
	})

	printer := infra.TerminalPrinter{}

	for _, cat := range categories {
		format := fmt.Sprintf("%%-%ds", nameLength)
		printer.Green(format, cat).
			Bold(" %s", infra.FormatDurationBracketed(durations[cat])).
			Newline()
	}

	return printer.String(), nil
}
