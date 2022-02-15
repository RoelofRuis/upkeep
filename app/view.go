package app

import (
	"fmt"
	"strings"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

func (r Repository) HandleViewWeek(args []string) (string, error) {
	date := model.NewDate(time.Now()).PreviousMonday()

	upkeep, err := r.Upkeep.Get()
	if err != nil {
		return "", err
	}

	sheets := make([]model.Timesheet, 5)
	for i, day := range date.IterateNext(5) {
		sheet, err := r.Timesheets.GetForDate(day)
		if err != nil {
			return "", err
		}
		sheets[i] = sheet
	}

	return ViewWeek(upkeep, sheets), nil
}

func ViewWeek(upkeep model.Upkeep, sheets []model.Timesheet) string {
	printer := infra.TerminalPrinter{}

	weekDur := time.Duration(0)
	weekQuotum := model.NewDuration()
	for _, daySheet := range sheets {
		blocks := upkeep.DiscountTimeBlocks(daySheet, time.Now())
		dayDur := blocks.TotalDuration()
		weekDur += dayDur

		dayQuotum := upkeep.GetTimesheetQuotum(daySheet)
		weekQuotum = weekQuotum.Add(dayQuotum)

		printer.Print("%s ", daySheet.Date.Format("Mon 02 Jan 2006"))

		if dayDur == 0 && !dayQuotum.IsDefined() {
			printer.Newline()
			continue
		}

		if !dayQuotum.IsDefined() {
			printer.Bold("[%s]", infra.FormatDuration(dayDur))
		} else {
			printer.Bold("[%s]", infra.FormatDuration(dayDur)).
				Print(" / [%s] ", infra.FormatDuration(dayQuotum.Get()))
		}

		printer.Green("%s", strings.Join(daySheet.GetCategoryNames(), " ")).Newline()
	}

	weekPerc := (float64(weekDur) / float64(weekQuotum.Get())) * 100

	printer.Bold("                [%s]", infra.FormatDuration(weekDur)).
		Print(" / [%s] (%0.1f%%)", infra.FormatDuration(weekQuotum.Get()), weekPerc)

	return printer.String()
}

func (r Repository) HandleViewSheet(args []string) (string, error) {
	date := model.NewDate(time.Now())
	if len(args) > 0 {
		switch args[0] {
		case "today":
			break
		case "yesterday":
			date = date.ShiftDay(-1)
			break
		default:
			parsedDate, err := model.NewDateFromString(args[0])
			if err != nil {
				return "", fmt.Errorf("invalid date value '%s'", args[0])
			}
			date = parsedDate
			break
		}
	}

	upkeep, err := r.Upkeep.Get()
	if err != nil {
		return "", err
	}
	timesheet, err := r.Timesheets.GetForDate(date)
	if err != nil {
		return "", err
	}

	return ViewSheet(upkeep, timesheet), nil
}

func ViewSheet(upkeep model.Upkeep, timesheet model.Timesheet) string {
	printer := infra.TerminalPrinter{}

	printer.Bold("@ %s", timesheet.Date.Format("Monday 02 Jan 2006")).Newline()
	printer.BGGreen("%s", strings.Join(upkeep.SelectedCategories, " | ")).Newline()

	blocks := upkeep.DiscountTimeBlocks(timesheet, time.Now())

	for _, block := range blocks {
		if block.Block.Id == -1 {
			printer.White(">> ")
		} else {
			printer.White("%2d ", block.Block.Id)
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
					printer.Yellow(" [%s] ", infra.FormatDuration(block.DiscountedDuration))
				} else {
					printer.Print("        ")
				}
			} else {
				printer.Bold(" [%s] ", infra.FormatDuration(block.Block.BaseDuration()))
			}
		} else {
			printer.Red("[%s -   ?  ]        ",
				block.Block.WithTime.Start.Format(model.LayoutHour),
			)
		}

		printer.Green("%s", block.Block.Category)

		printer.Newline()
	}

	quotum := upkeep.GetTimesheetQuotum(timesheet)
	totalDuration := blocks.TotalDuration()

	if !quotum.IsDefined() {
		printer.Print("                   ").
			Bold("[%s]", infra.FormatDuration(totalDuration)).
			Newline()
	} else {
		perc := (float64(totalDuration) / float64(quotum.Get())) * 100

		printer.Print("                   ").
			Bold("[%s]", infra.FormatDuration(totalDuration)).
			Print(" / [%s] (%0.1f%%)", infra.FormatDuration(quotum.Get()), perc).
			Newline()
	}

	return printer.String()
}

func (r Repository) HandleViewCategories(args []string) (string, error) {
	date := model.NewDate(time.Now())

	upkeep, err := r.Upkeep.Get()
	if err != nil {
		return "", err
	}
	timesheet, err := r.Timesheets.GetForDate(date)
	if err != nil {
		return "", err
	}

	return ViewCategories(upkeep, []model.Timesheet{timesheet}), nil
}

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
