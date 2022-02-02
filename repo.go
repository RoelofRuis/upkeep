package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"
	"timesheet/model"
)

type TimesheetRepository struct {
	fileIO FileIO
}

func (r *TimesheetRepository) GetForDay(t time.Time) (*model.Timesheet, error) {
	day := t.Format("2006-01-02")
	sheet := model.NewTimesheet(day)
	f, err := r.fileIO.OpenForDay(day)
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

	// Load timesheet fields
	breakActive, err := strconv.ParseBool(rows[0][0])
	if err != nil {
		return nil, err
	}
	lastStart, err := model.NewMomentFromString(rows[0][1])
	if err != nil {
		return nil, err
	}

	sheet.Break = breakActive
	sheet.Tags = model.NewTagSetFromString(rows[0][2])
	sheet.LastStart = lastStart

	// load blocks
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
		tags := model.NewTagSetFromString(row[2])
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
	f, err := r.fileIO.CreateForDay(t.Day)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	var rows [][]string

	// save timesheet fields
	rows = append(rows, []string{
		strconv.FormatBool(t.Break),
		t.LastStart.String(),
		t.Tags.String(),
	})

	// save blocks
	for _, block := range t.Blocks {
		rows = append(rows, []string{
			block.Start.String(),
			block.End.String(),
			block.Tags.String(),
		})
	}

	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}

	return nil
}
