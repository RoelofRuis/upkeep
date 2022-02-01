package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"timesheet/model"
)

type TimesheetRepository struct {
	path string
}

func (r *TimesheetRepository) GetForDay(t time.Time) (*model.Timesheet, error) {
	day := t.Format("2006-01-02")
	sheet := model.NewTimesheet(day)
	f, err := os.Open(fmt.Sprintf("%s/sheet_%s.csv", r.path, day))
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return sheet, nil
		default:
			return nil, err
		}
	}

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, errors.New("invalid timesheet file")
	}

	lastStart, err := model.NewMomentFromString(rows[0][1])
	if err != nil {
		return nil, err
	}

	sheet.Day = day
	sheet.LastStart = lastStart

	var blocks []model.TimeBlock

	for _, row := range rows[1:] {
		start, err := model.NewMomentFromString(row[0])
		if err != nil {
			return nil, err
		}
		end, err := model.NewMomentFromString(row[1])
		if err != nil {
			return nil, err
		}
		tags := []string{}
		if row[2] != "" {
			tags = strings.Split(row[2], ",")
		}
		blocks = append(blocks, model.TimeBlock{
			Start: start,
			End:   end,
			Tags:  tags,
		})
	}

	sheet.Blocks = blocks

	return sheet, nil
}

func (r *TimesheetRepository) Insert(t *model.Timesheet) error {
	f, err := os.Create(fmt.Sprintf("%s/sheet_%s.csv", r.path, t.Day))
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	var rows [][]string

	rows = append(rows, []string{
		t.Day,
		t.LastStart.String(),
	})

	for _, block := range t.Blocks {
		rows = append(rows, []string{
			block.Start.String(),
			block.End.String(),
			strings.Join(block.Tags, ","),
		})
	}

	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}

	return nil
}
