package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

func (a *Dependencies) Export() func(req *Request) (string, error) {
	return func(req *Request) (string, error) {
		groupCategories := GroupCategories(req.Params)

		categoryTotals := make(map[string]time.Duration)
		allDays := make(map[model.Date]map[string]time.Duration)

		for _, sheet := range req.Timesheets {
			dateDurs := make(map[string]time.Duration)

			blocks := req.Upkeep.DiscountTimeBlocks(*sheet, req.Clock.Now())
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
		for name, dur := range categoryTotals {
			if dur == 0 {
				continue
			}

			categoryNames = append(categoryNames, name)
		}

		sort.Strings(categoryNames)

		var records [][]string

		format := req.Params.GetNamed("f", "")
		if format == "excel" {
			records = append(records, []string{"sep=,"})
		}

		headers := []string{"DATE", "DAY"}
		for _, name := range categoryNames {
			headers = append(headers, name)
		}
		quotaSum := model.NewDuration()
		headers = append(headers, "TOTALS", "QUOTA")
		records = append(records, headers)

		for _, sheet := range req.Timesheets {
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

			sheetQuotum := req.Upkeep.GetTimesheetQuotum(*sheet)
			quotaSum = quotaSum.Add(sheetQuotum)

			record = append(record, infra.FormatDuration(sumDur))
			if sheet.Quotum.IsZero() && !sheet.AdjustedQuotum {
				record = append(record, "")
			} else {
				record = append(record, infra.FormatDuration(sheet.Quotum.Get()))
			}

			if sumDur > 0 {
				records = append(records, record)
			}
		}

		totals := []string{"", "TOTALS"}
		var sumDur = model.NewDuration()

		for _, name := range categoryNames {
			dur, _ := categoryTotals[name]
			sumDur = sumDur.AddDuration(dur)
			totals = append(totals, infra.FormatDuration(dur))
		}
		totals = append(totals, infra.FormatDuration(sumDur.Get()))
		totals = append(totals, infra.FormatDuration(quotaSum.Get()))

		quotaDiff := sumDur.Sub(quotaSum).Get()
		if quotaDiff < 0 {
			totals = append(totals, fmt.Sprintf("%s short", infra.FormatDuration(-quotaDiff)))
		} else {
			totals = append(totals, fmt.Sprintf("%s extra", infra.FormatDuration(quotaDiff)))
		}
		records = append(records, totals)

		percentages := []string{"", "PERCENTAGES"}
		for _, name := range categoryNames {
			dur, _ := categoryTotals[name]
			perc := float64(dur) / float64(sumDur.Get()) * 100
			percentages = append(percentages, fmt.Sprintf("%0.2f%%", perc))
		}
		records = append(records, percentages)

		// export records
		exportName := fmt.Sprintf("export_%s.csv", req.Clock.Now().Format("20060102_150405"))
		if err := a.IO.Export(exportName, records); err != nil {
			return "", err
		}

		return fmt.Sprintf("Wrote %s to current working directory", exportName), nil
	}
}
