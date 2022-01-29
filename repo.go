package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type TimesheetRepository struct {
	path string
}

func (r *TimesheetRepository) GetForDay(t time.Time) (*Timesheet, error) {
	day := t.Format("2006-01-02")
	f, err := os.Open(fmt.Sprintf("%s/sheet_%s.csv", r.path, day))
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return &Timesheet{Day: day}, nil
		default:
			return nil, err
		}
	}

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var blocks []TimeBlock

	for _, row := range rows {
		start, err := NewMomentFromString(row[0])
		if err != nil {
			return nil, err
		}
		end, err := NewMomentFromString(row[1])
		if err != nil {
			return nil, err
		}
		tags := []string{}
		if row[2] != "" {
			tags = strings.Split(row[2], ",")
		}
		blocks = append(blocks, TimeBlock{
			Start: start,
			End:   end,
			Tags:  tags,
		})
	}

	sheet := &Timesheet{
		Day:    day,
		Blocks: blocks,
	}

	return sheet, nil
}

func (r *TimesheetRepository) Insert(t *Timesheet) error {
	f, err := os.Create(fmt.Sprintf("%s/sheet_%s.csv", r.path, t.Day))
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	var rows [][]string
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
