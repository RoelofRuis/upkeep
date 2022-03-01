package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

func Export(io infra.FileIO) func(app *App) (string, error) {
	return func(app *App) (string, error) {
		groupCategories := GroupCategories(app.Params)

		categoryTotals := make(map[string]time.Duration)
		allDays := make(map[model.Date]map[string]time.Duration)

		for _, sheet := range app.Timesheets {
			dateDurs := make(map[string]time.Duration)

			blocks := app.Upkeep.DiscountTimeBlocks(*sheet, time.Now())
			for _, block := range blocks {
				category := block.Block.Category.GetName(groupCategories)
				dur, has := dateDurs[category]
				if !has {
					dur = time.Duration(0)
				}
				dateDurs[category] = dur + block.DiscountedDuration

				catDur, has := categoryTotals[category]
				if !has {
					catDur = time.Duration(0)
				}
				categoryTotals[category] = catDur + block.DiscountedDuration
			}

			allDays[sheet.Date] = dateDurs
		}

		var categoryNames []string
		for name := range categoryTotals {
			categoryNames = append(categoryNames, name)
		}

		sort.Strings(categoryNames)

		var records [][]string
		headers := []string{"DATE", "DAY"}
		for _, name := range categoryNames {
			headers = append(headers, name)
		}
		headers = append(headers, "TOTALS")
		records = append(records, headers)

		for _, sheet := range app.Timesheets {
			record := []string{
				sheet.Date.Format("2006-01-02"),
				sheet.Date.Format("Monday"),
			}
			dayCategories := allDays[sheet.Date]
			var sumDur = time.Duration(0)
			for _, name := range categoryNames {
				dur, has := dayCategories[name]
				if !has {
					dur = time.Duration(0)
				}
				sumDur += dur
				if dur == 0 {
					record = append(record, "")
				} else {
					record = append(record, infra.FormatDuration(dur))
				}
			}
			if sumDur > 0 {
				record = append(record, infra.FormatDuration(sumDur))
				records = append(records, record)
			}
		}

		totals := []string{"", "TOTALS"}
		var sumDur = time.Duration(0)

		for _, name := range categoryNames {
			dur, _ := categoryTotals[name]
			sumDur += dur
			totals = append(totals, infra.FormatDuration(dur))
		}
		totals = append(totals, infra.FormatDuration(sumDur))
		records = append(records, totals)

		percentages := []string{"", "PERCENTAGES"}
		for _, name := range categoryNames {
			dur, _ := categoryTotals[name]
			perc := float64(dur) / float64(sumDur) * 100
			percentages = append(percentages, fmt.Sprintf("%0.2f%%", perc))
		}
		records = append(records, percentages)

		// export records
		if err := io.Export("export.csv", records); err != nil {
			return "", err
		}

		return "Wrote export.csv", nil
	}
}
