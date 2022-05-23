package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
)

func ViewCategories(req *Request) (string, error) {
	var categories []string
	durations := make(map[string]time.Duration)
	nameLength := 0

	groupCategories := GroupCategories(req.Params)

	for _, sheet := range req.Timesheets {
		blocks := req.Upkeep.DiscountTimeBlocks(*sheet, req.Clock.Now())
		for _, block := range blocks {
			category := block.Block.Category.GetName(groupCategories)
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

	totalDur := time.Duration(0)
	for _, cat := range categories {
		totalDur += durations[cat]
	}

	sortType := req.Params.GetNamed("sort", "%")
	switch sortType {
	case "a":
		sort.Slice(categories, func(i, j int) bool {
			return categories[i] < categories[j]
		})

	default:
		sort.Slice(categories, func(i, j int) bool {
			return durations[categories[i]] > durations[categories[j]]
		})
	}

	printer := infra.TerminalPrinter{}

	for _, cat := range categories {
		format := fmt.Sprintf("%%-%ds", nameLength)
		printer.PrintC(infra.Green, format, cat).
			PrintC(infra.Bold, " %s", infra.FormatDurationBracketed(durations[cat])).
			Print(" %s", infra.FormatPercentage(durations[cat], totalDur)).
			Newline()
	}

	printer.Repeat(" ", nameLength).
		PrintC(infra.Bold, " %s", infra.FormatDurationBracketed(totalDur))

	return printer.String(), nil
}
